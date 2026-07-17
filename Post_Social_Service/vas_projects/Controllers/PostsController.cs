using Core.DTO.Response;
using Core.Features.Posts.Queries.GetFeedPosts;
using Core.Interfaces;
using MediatR;
using Microsoft.AspNetCore.Mvc;

namespace vas_projects.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class PostsController(ILogger<PostsController> _logger, ISender _sender) : ControllerBase
    {
        [HttpGet("feed")]
        public async Task<IActionResult> GetFeed([FromQuery] Guid userId, CancellationToken cancellation)
        {
            var query = new GetFeedQuery(userId);
            var posts = await _sender.Send(query, cancellation);
            _logger.LogInformation($"Service feeded post for User with '{userId}' Id");
            return Ok(posts);
        }
    }
}
