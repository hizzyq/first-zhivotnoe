using System;
using System.Collections.Generic;
using System.Text;
using MediatR; 

namespace Core.Features.Posts.Commands.CreatePost
{
    public record CreatePostCommand(Guid UserId, string Title, string Description, string ImagePath) : IRequest<Guid>;
}
