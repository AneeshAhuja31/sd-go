-- POSTS TABLE

CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    content TEXT NOT NULL,
    views INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- INDEX 

CREATE INDEX IF NOT EXISTS idx_posts_email
ON posts(email);


INSERT INTO posts (email, content, views) VALUES
-- aneesh
('aneesh@example.com', 'Hello world!', 10),
('aneesh@example.com', 'Building microservices in Go', 34),
('aneesh@example.com', 'Circuit breakers explained', 21),

-- rahul
('rahul@example.com', 'Golang is fast', 18),
('rahul@example.com', 'Why I love channels', 27),
('rahul@example.com', 'Writing clean APIs', 12),

-- neha
('neha@example.com', 'Full-stack journey begins', 9),
('neha@example.com', 'React + Go combo', 22),
('neha@example.com', 'System design notes', 31),

-- arjun
('arjun@example.com', 'Distributed systems are hard', 45),
('arjun@example.com', 'CAP theorem demystified', 60),
('arjun@example.com', 'Eventual consistency in practice', 39),

-- kavya
('kavya@example.com', 'Getting started with ML', 14),
('kavya@example.com', 'Feature engineering basics', 26),
('kavya@example.com', 'Model deployment tips', 33),

-- rohit
('rohit@example.com', 'Cloud basics', 11),
('rohit@example.com', 'Docker for beginners', 24),
('rohit@example.com', 'Kubernetes concepts', 41),

-- sneha
('sneha@example.com', 'Frontend performance tips', 19),
('sneha@example.com', 'CSS architecture', 23),
('sneha@example.com', 'Design systems 101', 29),

-- vikram
('vikram@example.com', 'API versioning strategies', 17),
('vikram@example.com', 'REST vs gRPC', 36),
('vikram@example.com', 'Error handling patterns', 28),

-- pooja
('pooja@example.com', 'Data pipelines overview', 20),
('pooja@example.com', 'Batch vs streaming', 32),
('pooja@example.com', 'SQL optimization tricks', 44),

-- aman
('aman@example.com', 'Systems programming basics', 13),
('aman@example.com', 'Memory management in Go', 25),
('aman@example.com', 'Performance tuning tips', 37),

-- ishan
('ishan@example.com', 'Latency budgeting', 34),
('ishan@example.com', 'High-performance services', 40),
('ishan@example.com', 'Load testing strategies', 48),

-- kriti
('kriti@example.com', 'Shipping fast safely', 21),
('kriti@example.com', 'Feature prioritization', 27),
('kriti@example.com', 'Building products users love', 16);
