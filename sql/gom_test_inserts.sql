-- Some test data to insert into the GOM database for testing purposes
call usp_insert_project('GOMTest', 'Test Project for GOM', Now(), CURRENT_USER);
call usp_insert_protocol('FTP', 'Standard FTP Connection', Now(), CURRENT_USER);
call usp_insert_protocol('SFTP', 'Secure FTP Connection', Now(), CURRENT_USER);
call usp_insert_protocol('S3', 'AWS Simple Storage Service', Now(), CURRENT_USER);
call usp_insert_auth_detail('Testing Values', 'gom_user', 'gom_user_password', '$HOME/.ssh/_id_ed25519', '', '', Now(), CURRENT_USER);
call usp_insert_connection(2, 1, Now(), CURRENT_USER);
call usp_insert_schedule('GOM Test Schedule', Now(), CURRENT_USER);
call usp_insert_action('/app/data/test', '/test/path/foo', 'localhost', 1, 1, Now(), CURRENT_USER);
call usp_insert_schedule_action(1, 1);