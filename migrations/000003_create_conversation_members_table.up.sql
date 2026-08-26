CREATE TABLE conversation_members (
                                      id BIGSERIAL PRIMARY KEY,
                                      conversation_id BIGINT NOT NULL,
                                      user_id BIGINT NOT NULL,

                                      CONSTRAINT fk_conversation_members_conversation
                                          FOREIGN KEY (conversation_id)
                                              REFERENCES conversations(id)
                                              ON DELETE CASCADE,

                                      CONSTRAINT fk_conversation_members_user
                                          FOREIGN KEY (user_id)
                                              REFERENCES users(id)
                                              ON DELETE CASCADE,

                                      CONSTRAINT uq_conversation_member
                                          UNIQUE (conversation_id, user_id)
);

CREATE INDEX idx_conversation_members_user_id
    ON conversation_members(user_id);

CREATE INDEX idx_conversation_members_conversation_id
    ON conversation_members(conversation_id);