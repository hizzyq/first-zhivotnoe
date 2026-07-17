using Date.Entities;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace Date.Configurations
{
    public class PostConfiguration : IEntityTypeConfiguration<PostsEntity>
    {
        public void Configure(EntityTypeBuilder<PostsEntity> builder)
        {
            builder.ToTable("Posts");

            builder.HasKey(p => p.PostId);

            builder.Property(p => p.Title)
                .IsRequired()
                .HasMaxLength(250);

            builder.Property(p => p.Description)
                .IsRequired();

            builder.HasIndex(p => p.UserId);
        }
    }
}
