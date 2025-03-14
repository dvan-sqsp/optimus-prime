CREATE TABLE pull_requests (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   github_id INTEGER NOT NULL,
   repository_name VARCHAR(255) NOT NULL,
   repository_owner VARCHAR(255) NOT NULL,
   title VARCHAR(255) NOT NULL,
   description TEXT,
   state VARCHAR(20) NOT NULL, -- open, closed, merged
   author_username VARCHAR(255) NOT NULL,
   author_id INTEGER NOT NULL,
   created_at TIMESTAMP WITH TIME ZONE NOT NULL,
   updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
   closed_at TIMESTAMP WITH TIME ZONE,
   merged_at TIMESTAMP WITH TIME ZONE,
   source_branch VARCHAR(255) NOT NULL,
   target_branch VARCHAR(255) NOT NULL,
   additions_count INTEGER,
   deletions_count INTEGER,
   changed_files_count INTEGER,
   review_status VARCHAR(20), -- pending, approved, changes_requested
   url VARCHAR(255) NOT NULL,
   UNIQUE(repository_owner, repository_name, github_id)
);