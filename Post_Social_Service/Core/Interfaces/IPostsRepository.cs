using Core.Models;

namespace Core.Interfaces
{
    public interface IPostsRepository
    {
        Task CreatePost(Post _post, CancellationToken cancellation);
    }
}