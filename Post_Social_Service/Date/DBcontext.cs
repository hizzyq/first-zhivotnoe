using Date.Configurations;
using Date.Entities;
using Microsoft.EntityFrameworkCore;

namespace Date
{
    public class AppDbContext : DbContext
    {
        public AppDbContext(DbContextOptions<AppDbContext> options) : base(options) { }

        public DbSet<PostsEntity> Posts {  get; set; }

        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            base.OnModelCreating(modelBuilder);

            modelBuilder.ApplyConfiguration<PostsEntity>(new PostConfiguration());
        }
    }
}
