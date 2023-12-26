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


CREATE OR REPLACE FUNCTION create_worker(
    p_login VARCHAR(255),
    p_surname VARCHAR(64),
    p_name VARCHAR(64),
    p_patronymic VARCHAR(64),
    p_phone VARCHAR(15),
    p_position UUID,
    p_department UUID
) RETURNS VOID AS $$
DECLARE
    v_id UUID;
BEGIN
    INSERT INTO staff(surname, name, patronymic, phone_number, position, department) VALUES
    (p_surname, p_name, p_patronymic, p_phone, p_position, p_department)
    RETURNING id INTO v_id;

    INSERT INTO accounts(login, employee)
    VALUES (p_login, v_id);
END;
$$ LANGUAGE PLPGSQL;

CREATE OR REPLACE FUNCTION activate_account(p_login VARCHAR(255), p_password VARCHAR(1024))
RETURNS UUID
AS $$
DECLARE
    v_employee_id UUID;
BEGIN
    SELECT employee INTO v_employee_id
    FROM accounts
    WHERE login = p_login
    FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'User with login % not found', p_login;
    END IF;

    IF EXISTS (SELECT 1 FROM accounts WHERE login = p_login AND activated) THEN
        RAISE EXCEPTION 'User with login % is already activated', p_login;
    END IF;

    UPDATE accounts
    SET password = p_password, activated = TRUE
    WHERE login = p_login;

    RETURN v_employee_id;
END;
$$ LANGUAGE plpgsql;
