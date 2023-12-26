#!/usr/bin/bash

# Since the Postgres currently unable to use ENV Variables inside .sql scripts.
# This file contains definitions user roles, values of which are taken from ENV variables

psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} <<-END
    CREATE ROLE ${POSTGRES_ADMIN_NAME} WITH
    LOGIN
    PASSWORD '${POSTGRES_ADMIN_PWD}';

    GRANT ALL PRIVILEGES ON TABLE positions TO ${POSTGRES_ADMIN_NAME};
    GRANT ALL PRIVILEGES ON TABLE departments TO ${POSTGRES_ADMIN_NAME};
    GRANT ALL PRIVILEGES ON TABLE professional_developments TO ${POSTGRES_ADMIN_NAME};
    GRANT ALL PRIVILEGES ON TABLE incidents TO ${POSTGRES_ADMIN_NAME};
    GRANT ALL PRIVILEGES ON TABLE accounts TO ${POSTGRES_ADMIN_NAME};
    GRANT ALL PRIVILEGES ON TABLE staff TO ${POSTGRES_ADMIN_NAME};
    

    CREATE ROLE ${POSTGRES_HEAD_NAME} WITH
    LOGIN
    PASSWORD '${POSTGRES_HEAD_PWD}';

    GRANT SELECT ON TABLE positions TO ${POSTGRES_HEAD_NAME}; -- read positions
    GRANT SELECT ON TABLE departments TO ${POSTGRES_HEAD_NAME}; -- read departments
    GRANT SELECT, INSERT, UPDATE ON TABLE professional_developments TO ${POSTGRES_HEAD_NAME}; -- read and add/update developments for staff
    GRANT SELECT, UPDATE ON TABLE accounts TO ${POSTGRES_HEAD_NAME}; -- read and update it's own account
    GRANT SELECT, UPDATE ON staff TO ${POSTGRES_HEAD_NAME}; -- read and update staff 
    GRANT SELECT, INSERT, UPDATE ON TABLE incidents TO ${POSTGRES_STAFF_NAME}; -- add and update department's incidents


    CREATE ROLE ${POSTGRES_STAFF_NAME} WITH
    LOGIN 
    PASSWORD '${POSTGRES_STAFF_PWD}';

    GRANT SELECT ON TABLE positions TO ${POSTGRES_STAFF_NAME}; -- read positions
    GRANT SELECT ON TABLE departments TO ${POSTGRES_STAFF_NAME}; -- read departments
    GRANT SELECT ON TABLE professional_developments TO ${POSTGRES_STAFF_NAME}; -- read it's own proff-dev courses
    GRANT SELECT, UPDATE ON TABLE accounts TO ${POSTGRES_STAFF_NAME}; -- read and update it's own account(password f.e.)
    GRANT SELECT, INSERT, UPDATE ON TABLE incidents TO ${POSTGRES_STAFF_NAME}; -- add and update it's own incidents
END
