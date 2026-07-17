using Confluent.Kafka;
using Core.Features.Posts.Commands.CreatePost;
using MediatR;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using System.Text.Json;

namespace vas_projects.Infrastructure
{
    public class KafkaConsumerService : BackgroundService
    {
        private readonly ILogger<KafkaConsumerService> _logger;
        private IConsumer<Ignore, string> _consumer;
        private readonly string _topicName;
        private readonly IServiceProvider _serviceProvider;
        private readonly IConfiguration _configuration;

        public KafkaConsumerService(ILogger<KafkaConsumerService> logger, IConfiguration config, IServiceProvider provider)
        {
            _logger = logger;
            _serviceProvider = provider;
            _topicName = config["KafkaSettings:Topic"] ?? "__consumer_offsets";
            _configuration = config;
        }

        protected override Task ExecuteAsync(CancellationToken stoppingToken)
        {
            return Task.Run(() => StartConsumerLoop(stoppingToken), stoppingToken);
        }

        private async Task StartConsumerLoop(CancellationToken stoppingToken) 
        {

            try
            {
                var bootstrapServise = _configuration["KafkaSettings:BootstrapServer"];
                var groupId = _configuration["KafkaSettings:GroupId"];

                var cfg = new ConsumerConfig
                {
                    BootstrapServers = bootstrapServise,
                    GroupId = groupId,
                    AutoOffsetReset = AutoOffsetReset.Earliest,
                    EnableAutoCommit = false
                };

                _consumer = new ConsumerBuilder<Ignore, string>(cfg).Build();
            }
            catch (Exception ex) 
            {
                _logger.LogCritical(ex, "Cannot initialize Kafka Consumer");
                return;
            }
            _consumer.Subscribe(_topicName);
            _logger.LogInformation($"Subscribe on topic {_topicName}. Wait a message... ");

            try
            {
                while (!stoppingToken.IsCancellationRequested)
                {
                    try
                    {
                        var consumerResult = _consumer.Consume(stoppingToken);

                        if (consumerResult != null)
                        {
                            var rawJson = consumerResult.Message.Value;
                            _logger.LogInformation($"Got a message: {consumerResult.Message.Value}");

                            var command = JsonSerializer.Deserialize<CreatePostCommand>(rawJson, new JsonSerializerOptions
                            {
                                PropertyNameCaseInsensitive = true
                            });

                            if (command == null)
                            {
                                _logger.LogError($"Cannot deserializes message: {rawJson}");
                                continue;
                            }

                            using (var scope = _serviceProvider.CreateScope())
                            {
                                var mediator = scope.ServiceProvider.GetRequiredService<IMediator>();

                                _logger.LogInformation($"Sending a command for creating a post: {command.Title}");
                                await mediator.Send(command, stoppingToken);
                            }

                            _consumer.Commit(consumerResult);
                            _logger.LogInformation("The message succesful processed and comited!");
                        }
                    }

                    catch (ConsumeException e)
                    {
                        _logger.LogError($"Error during got a message: {e.Message}");
                    }
                    catch (Exception e)
                    {
                        _logger.LogError("Error with Kafka service");
                    }
                }
            }
            catch (OperationCanceledException)
            {
                _logger.LogError($"Stopped listing Kafka, application stopped work");
            }
            finally
            {
                _consumer.Close();
            }
        }

        public override void Dispose() 
        {
            _consumer.Dispose();
            base.Dispose();
        }
    }
}
