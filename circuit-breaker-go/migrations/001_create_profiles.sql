-- PROFILES TABLE

CREATE TABLE IF NOT EXISTS profiles (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL,
    dob DATE,
    bio TEXT,
    hobbies TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

INSERT INTO profiles (email, username, dob, bio) VALUES
('aneesh@example.com',  'aneesh',  '1999-01-01', 'Backend developer'),
('rahul@example.com',   'rahul',   '1998-05-12', 'Golang enthusiast'),
('neha@example.com',    'neha',    '2000-03-21', 'Full-stack engineer'),
('arjun@example.com',   'arjun',   '1997-11-09', 'Distributed systems fan'),
('kavya@example.com',   'kavya',   '1999-07-18', 'Machine learning explorer'),
('rohit@example.com',   'rohit',   '1996-02-14', 'Cloud & DevOps'),
('sneha@example.com',   'sneha',   '2001-09-30', 'Frontend specialist'),
('vikram@example.com',  'vikram',  '1998-12-05', 'API designer'),
('pooja@example.com',   'pooja',   '1997-04-27', 'Data engineer'),
('aman@example.com',    'aman',    '1995-08-16', 'Systems programmer'),
('ishan@example.com',   'ishan',   '1999-10-22', 'High-performance services'),
('kriti@example.com',   'kriti',   '2000-06-11', 'Product-focused engineer'),
('aditya@example.com',  'aditya',  '1996-01-19', 'Microservices architect'),
('nisha@example.com',   'nisha',   '1998-03-03', 'Reliability engineer'),
('karan@example.com',   'karan',   '1999-12-29', 'AI + backend'),
('megha@example.com',   'megha',   '2001-05-07', 'Learning full stack'),
('saurabh@example.com', 'saurabh', '1997-09-14', 'Networking nerd'),
('rhea@example.com',    'rhea',    '1998-02-25', 'UX-minded developer'),
('manish@example.com',  'manish',  '1996-07-01', 'Scalability geek'),
('tanvi@example.com',   'tanvi',   '2000-11-20', 'Clean code advocate'),
('pranav@example.com',  'pranav',  '1999-04-04', 'Backend intern');

UPDATE profiles
SET hobbies = CASE email
WHEN 'aneesh@example.com'  THEN ARRAY['backend','golang','microservices','system_design']
WHEN 'rahul@example.com'   THEN ARRAY['golang','backend','concurrency','performance']
WHEN 'neha@example.com'    THEN ARRAY['frontend','backend','databases','fullstack']
WHEN 'arjun@example.com'   THEN ARRAY['distributed_systems','backend','scalability']
WHEN 'kavya@example.com'   THEN ARRAY['ai_ml','machine_learning','backend']
WHEN 'rohit@example.com'   THEN ARRAY['cloud','devops','docker','kubernetes']
WHEN 'sneha@example.com'   THEN ARRAY['frontend','ui','design_systems']
WHEN 'vikram@example.com'  THEN ARRAY['apis','backend','system_design']
WHEN 'pooja@example.com'   THEN ARRAY['data_engineering','sql','databases']
WHEN 'aman@example.com'    THEN ARRAY['systems','backend','performance']
WHEN 'ishan@example.com'   THEN ARRAY['performance','scalability','distributed_systems']
WHEN 'kriti@example.com'   THEN ARRAY['product_engineering','clean_code','backend']
WHEN 'aditya@example.com'  THEN ARRAY['microservices','scalability','cloud']
WHEN 'nisha@example.com'   THEN ARRAY['reliability','sre','monitoring']
WHEN 'karan@example.com'   THEN ARRAY['ai_ml','backend','model_serving']
WHEN 'megha@example.com'   THEN ARRAY['fullstack','learning','frontend']
WHEN 'saurabh@example.com' THEN ARRAY['networking','http','performance']
WHEN 'rhea@example.com'    THEN ARRAY['ux','design','frontend']
WHEN 'manish@example.com'  THEN ARRAY['scalability','caching','distributed_systems']
WHEN 'tanvi@example.com'   THEN ARRAY['clean_code','refactoring','golang']
WHEN 'pranav@example.com'  THEN ARRAY['backend','learning','sql']
ELSE hobbies
END;

-- INDEX FOR SIMILARITY QUERIES

CREATE INDEX IF NOT EXISTS idx_profiles_hobbies
ON profiles USING GIN (hobbies);
