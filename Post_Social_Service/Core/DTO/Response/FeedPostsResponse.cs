using Core.Models;
using System;
using System.Collections.Generic;
using System.Text;

namespace Core.DTO.Response
{
    public record FeedPostsResponse(Guid PostId, Guid UserId,
        string Title, string Description, string ImagePath,
        int CountLikes, PostStatus Status, DateTime CreatedAt);
}
