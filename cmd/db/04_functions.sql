CREATE OR REPLACE FUNCTION create_superuser(
    p_login VARCHAR(255),
    p_password VARCHAR(1024)
)
RETURNS UUID AS $$
DECLARE
    v_id UUID;
    v_department UUID;
    v_position UUID;
BEGIN
    v_id := NULL;

    SELECT id INTO v_department
    FROM departments
    WHERE name = 'Labor-Protection-System';

    SELECT id INTO v_position
    FROM positions
    WHERE level = 'admin';

    IF v_position IS NULL OR v_department IS NULL THEN
        RAISE EXCEPTION 'The department and/or the position for superuser do not exist';
    END IF;

    INSERT INTO staff(surname, name, phone_number, position, department)
    VALUES ('admin surname', 'admin name', '0000000000000', v_position, v_department)
    RETURNING id INTO v_id;

    INSERT INTO accounts(login, password, activated, employee)
    VALUES (p_login, p_password, TRUE, v_id);

    RETURN v_id;
END;
$$ LANGUAGE PLPGSQL;
