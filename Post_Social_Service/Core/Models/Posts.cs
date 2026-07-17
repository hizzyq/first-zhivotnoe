using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Runtime.CompilerServices;
using System.Text;

namespace Core.Models
{
    public class Post
    {
        private Post(Guid postId, Guid userId, string title, string desc, string imagePath)
        {
            PostId = postId;
            UserId = userId;
            Title = title;
            Description = desc;
            ImagePath = imagePath;
            Status = PostStatus.Active;
            CountLikes = 0;
        }

        private Post(Guid postId, Guid userId,
            string title, string desc, string imagePath,
            int countLikes, PostStatus status, DateTime datetime)
        {
            PostId = postId;
            UserId = userId;
            Title = title;
            Description = desc;
            ImagePath = imagePath;
            Status = status;
            CountLikes = countLikes;
            CreatedAt = datetime;
        }
        public Guid PostId { get; } 

        [Required]
        public Guid UserId { get; } 

        [Required]
        [MaxLength(50)]
        public string Title { get; }

        [MaxLength(1000)]
        public string Description { get; } 

        public string ImagePath { get;}

        public int CountLikes { get; } 

        public PostStatus Status { get; }

        public DateTime CreatedAt { get; } = DateTime.UtcNow;

        public bool IsReadToDisplay() => Status == PostStatus.Active;

        public static Post Create(Guid userId, string title, string desc, string imagePath)
        {
            return new Post(Guid.NewGuid(), userId, title, desc, imagePath);
        }

        public static Post MapFromEntity(
            Guid postId, Guid userId,
            string title, string desc, string imagePath,
            int countLikes, int status, DateTime datetime)
        {
            return new Post(postId, userId, title, desc, imagePath, countLikes, (PostStatus)status, datetime);
        }
    }

    public enum PostStatus
    {
        Processing = 0,
        Active,    
        Archived    
    }
}
