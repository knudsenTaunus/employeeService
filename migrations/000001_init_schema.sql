-- +migrate Up

CREATE TABLE `employees` (`id` INTEGER PRIMARY KEY AUTO_INCREMENT UNIQUE ,`first_name` VARCHAR(64), `last_name` VARCHAR(64), `salary` INTEGER, `birthday` DATE, employee_number INTEGER UNIQUE, entry_date DATE);
CREATE TABLE `companycars` (`id` INTEGER PRIMARY KEY AUTO_INCREMENT UNIQUE ,`manufacturer` VARCHAR(64), `type` VARCHAR(64), `number_plate` TEXT, employee_number INTEGER, FOREIGN KEY (employee_number) REFERENCES employees(employee_number));

-- +migrate Down

DROP TABLE IF EXISTS companycars;
DROP TABLE IF EXISTS employees;