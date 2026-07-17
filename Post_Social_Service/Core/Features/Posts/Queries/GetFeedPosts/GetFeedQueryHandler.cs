using Core.DTO.Response;
using Dapper;
using MediatR;
using System;
using System.Collections.Generic;
using System.Data;
using System.Text;

namespace Core.Features.Posts.Queries.GetFeedPosts
{
    public class GetFeedQueryHandler : IRequestHandler<GetFeedQuery, IEnumerable<FeedPostsResponse>>
    {
        private readonly IDbConnection _dbConnection;

        public GetFeedQueryHandler(IDbConnection dbConnection)
        {
            _dbConnection = dbConnection;
        }

        public async Task<IEnumerable<FeedPostsResponse>> Handle(GetFeedQuery query, CancellationToken cancellation)
        {
            const string sql = @"
                SELECT * FROM ""Posts"" WHERE ""UserId"" = @CurrentUserId";

            var command = new CommandDefinition(
                sql,
                parameters : new { CurrentUserId = query.UserId },
                cancellationToken : cancellation
                );

            return await _dbConnection.QueryAsync<FeedPostsResponse>(command);
        }
    }
}
