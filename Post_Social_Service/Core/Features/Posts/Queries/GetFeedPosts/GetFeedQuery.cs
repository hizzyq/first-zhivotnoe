using Core.DTO.Response;
using MediatR;


namespace Core.Features.Posts.Queries.GetFeedPosts
{
    public record GetFeedQuery(Guid UserId) : IRequest<IEnumerable<FeedPostsResponse>>;
}
