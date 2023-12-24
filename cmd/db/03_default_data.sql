
INSERT INTO positions(title, max_per_department) VALUES
    ('Department Head', 1),
    ('Project Owner', 1),
    ('Project Manager', NULL),
    ('Team-Lead', NULL),
    ('Developer', NULL),
    ('Analytic', NULL),
    ('Salesman', NULL),
    ('System Administrator', NULL),
    ('Administrator', 1); -- this is the LPS administrator

INSERT INTO departments(name) VALUES
    ('Labor-Protection-System'),
    ('Search Team'),
    ('Marketing');
