-- セッション管理テーブル
CREATE TABLE IF NOT EXISTS sessions (
    session_id CHAR(36) PRIMARY KEY,
    current_scene VARCHAR(10) NOT NULL DEFAULT 'S0',
    s6_started_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    INDEX idx_sessions_expires_at (expires_at),
    INDEX idx_sessions_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- S4回答テーブル
CREATE TABLE IF NOT EXISTS session_answers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    session_id CHAR(36) NOT NULL,
    question_id INT NOT NULL,
    answer_text TEXT NOT NULL,
    answered_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_session_question (session_id, question_id),
    INDEX idx_session_answers_session_id (session_id),
    FOREIGN KEY (session_id) REFERENCES sessions(session_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- S6進捗管理テーブル
CREATE TABLE IF NOT EXISTS session_s6_progress (
    id INT AUTO_INCREMENT PRIMARY KEY,
    session_id CHAR(36) NOT NULL,
    place_id VARCHAR(50) NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    verified_by VARCHAR(20),
    quiz_id CHAR(36),
    answered BOOLEAN NOT NULL DEFAULT FALSE,
    correct BOOLEAN NOT NULL DEFAULT FALSE,
    verified_at DATETIME,
    UNIQUE KEY uk_session_place (session_id, place_id),
    INDEX idx_s6_progress_session_id (session_id),
    FOREIGN KEY (session_id) REFERENCES sessions(session_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- クイズテーブル
CREATE TABLE IF NOT EXISTS quiz_questions (
    quiz_id CHAR(36) PRIMARY KEY,
    session_id CHAR(36) NOT NULL,
    place_id VARCHAR(50) NOT NULL,
    question_text TEXT NOT NULL,
    option_1 TEXT NOT NULL,
    option_2 TEXT NOT NULL,
    option_3 TEXT NOT NULL,
    option_4 TEXT NOT NULL,
    answer_index INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_quiz_questions_session_id (session_id),
    INDEX idx_quiz_questions_place_id (place_id),
    FOREIGN KEY (session_id) REFERENCES sessions(session_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- プレイヤーメッセージテーブル
CREATE TABLE IF NOT EXISTS player_messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    session_id CHAR(36) NOT NULL,
    place_id VARCHAR(50) NOT NULL,
    message_text TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_player_messages_place_id (place_id),
    INDEX idx_player_messages_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
