using Microsoft.Extensions.Logging;
using System;
using System.Collections.Generic;
using System.Text;
using Core.Models;
using Date.Entities;
using Microsoft.EntityFrameworkCore;
using Core.Interfaces;

namespace Date.Repositories
{
    public class PostsRepository : IPostsRepository
    {
        private readonly ILogger<PostsRepository> _logger;
        private readonly AppDbContext _context;

        public PostsRepository(ILogger<PostsRepository> logger,
            AppDbContext context)
        {
            _logger = logger;
            _context = context;
        }

        public async Task CreatePost(Post _post, CancellationToken cancellation)
        {
            _context.Posts.Add(new PostsEntity
            {
                PostId = _post.PostId,
                UserId = _post.UserId,
                Title = _post.Title,
                Description = _post.Description,
                ImagePath = _post.ImagePath,
                CountLikes = _post.CountLikes,
                Status = (int)_post.Status,
                CreatedAt = _post.CreatedAt,
            });

            await _context.SaveChangesAsync(cancellation);
        }
    }
}
