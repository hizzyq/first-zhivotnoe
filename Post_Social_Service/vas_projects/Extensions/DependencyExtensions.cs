using vas_projects.Infrastructure;
using Core.Interfaces;
using Date.Repositories;
using Npgsql;
using System.Data;

namespace vas_projects.Extensions
{
    public static class DependencyExtensions
    {
        public static IServiceCollection AddInfrastructureServices(this IServiceCollection services, string connectionString) 
        {
            services.AddScoped<IPostsRepository, PostsRepository>();
            services.AddScoped<IDbConnection>(sp => new NpgsqlConnection(connectionString));
            services.AddHostedService<KafkaConsumerService>();
            services.AddProblemDetails();

            return services;
        }
    }
}
