insert into account_type (name, is_active, description) values ('guest', true, 'an account without permit');
insert ignore into account (uuid, first_name, last_name, display_name, phone_number, password, email, address, type_id, role_id ,is_active, expire_at, description)
values ('8d816cfe-e263-4524-ae8d-a35785ed50ec','mohammad','boriaei','mohammad_admin','09100636393','123456','mohammad.boriaei@gmail.com','',(select id from account_type where name='guest'),(select id from role where name='administrator') ,true,'2030-01-01 15:15:15','this user can access to all resource as super admin.');
insert ignore into role (name, is_active,description) values ('administrator', 1,'this role is for service genesis');
insert ignore into role (name, is_active,description) values ('guest', 1,'this role is for guest user');
insert ignore into account_role(account_id, role_id, is_active,description) values ((select id from account where display_name='mohammad_admin'),(select id from role where name='administrator'), 1,'assign admin user to its roles');
