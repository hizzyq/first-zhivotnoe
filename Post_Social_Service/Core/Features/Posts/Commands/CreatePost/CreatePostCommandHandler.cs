using Core.Interfaces;
using MediatR;
using Core.Models;

namespace Core.Features.Posts.Commands.CreatePost
{
    class CreatePostCommandHandler : IRequestHandler<CreatePostCommand, Guid>
    {
        private readonly IPostsRepository _postRepos;

        public CreatePostCommandHandler(IPostsRepository postRepos) 
        {
            _postRepos = postRepos;
        }

        public async Task<Guid> Handle(CreatePostCommand command, CancellationToken cancellationToken)
        {
            Post post = Post.Create(
                command.UserId,
                command.Title,
                command.Description,
                command.ImagePath);

            await _postRepos.CreatePost(post, cancellationToken);

            return post.PostId;
        }
    }
}
