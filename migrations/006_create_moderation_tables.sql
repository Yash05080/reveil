CREATE TABLE IF NOT EXISTS public.moderation_flags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL REFERENCES public.community_posts(id) ON DELETE CASCADE,
    comment_id UUID, -- Nullable, for future use
    flag_reason TEXT NOT NULL,
    checked_by TEXT NOT NULL,
    action_taken TEXT NOT NULL,
    severity_level INTEGER NOT NULL DEFAULT 1,
    confidence_score FLOAT NOT NULL DEFAULT 1.0,
    notified_moderator BOOLEAN NOT NULL DEFAULT FALSE,
    flagged_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_moderation_flags_post_id ON public.moderation_flags(post_id);
CREATE INDEX IF NOT EXISTS idx_moderation_flags_reason ON public.moderation_flags(flag_reason);
