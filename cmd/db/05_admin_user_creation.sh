#!/usr/bin/bash

# This file contains an insert for admin role into db.

psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} <<-END
    SELECT create_superuser('administrator', 'password');
END
