USE master;

DROP DATABASE IF EXISTS Fortisure_Dev;

CREATE DATABASE Fortisure_Dev;

USE Fortisure_Dev;

  GO
CREATE SCHEMA Accounts
  GO
CREATE SCHEMA Company
  GO
CREATE SCHEMA Consulting
  GO
CREATE SCHEMA People
  GO
CREATE SCHEMA Sales
  GO
CREATE SCHEMA School
  GO
CREATE SCHEMA Security
  GO

CREATE TABLE dbo.access_keys
(
    kid char(27) PRIMARY KEY NOT NULL,
    x5c varchar(1500) NOT NULL,
    timestamp datetime DEFAULT getdate()
);

CREATE TABLE dbo.casbin_rule
(
    p_type nvarchar(100),
    v0 nvarchar(100),
    v1 nvarchar(100),
    v2 nvarchar(100),
    v3 nvarchar(100),
    v4 nvarchar(100),
    v5 nvarchar(100)
);

CREATE TABLE Company.Departments
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    name varchar(50) NOT NULL,
    description int
);


CREATE UNIQUE INDEX Departments_name_uindex ON Company.Departments (name);

CREATE TABLE People.Users
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    first_name varchar(50) NOT NULL,
    last_name varchar(50),
    username varchar(100) NOT NULL,
    ad_object_id char(36)
);

CREATE UNIQUE INDEX UQ_Usernames ON People.Users (username);
CREATE INDEX IX_Ad_Obj_ID ON People.Users (ad_object_id);

CREATE TABLE Company.Employee_Info
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    user_id int NOT NULL,
    title varchar(50) NOT NULL,
    department_id int,
    manager_id int,
    start_date date DEFAULT CONVERT([date],getdate()) NOT NULL,
    CONSTRAINT Employee_Info_Users_id_fk FOREIGN KEY (user_id) REFERENCES People.Users (id),
    CONSTRAINT Employee_Info_Departments_id_fk FOREIGN KEY (department_id) REFERENCES Company.Departments (id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT Employee_Info_Users_id_fk_2 FOREIGN KEY (manager_id) REFERENCES People.Users (id) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE Company.Positions
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    title varchar(100) NOT NULL,
    description varchar(2000)
);

CREATE TABLE People.User_Info
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    user_id int NOT NULL,
    preferred_name varchar(50),
    gender varchar(10),
    date_of_birth date,
    email varchar(100) NOT NULL,
    phone_primary varchar(20),
    phone_secondary varchar(20),
    address1 varchar(50),
    address2 varchar(50),
    city varchar(50),
    state char(2),
    zip varchar(10),
    bio text,
    CONSTRAINT FK_User_Info_Users FOREIGN KEY (user_id) REFERENCES People.Users (id)
);


CREATE TABLE Sales.Orders
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    user_id int NOT NULL,
    status varchar(50) NOT NULL,
    notes varchar(2000),
    created_on datetime DEFAULT getdate() NOT NULL,
    updated_on datetime,
    CONSTRAINT FK_Orders_Users FOREIGN KEY (user_id) REFERENCES People.Users (id)
);

  CREATE TABLE Sales.Products
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    price decimal(18,2) NOT NULL,
    name varchar(100) NOT NULL,
    category varchar(50) NOT NULL,
    description varchar(2000)
);
CREATE UNIQUE INDEX Products_name_uindex ON Sales.Products (name);


 CREATE TABLE Sales.Order_Line_Items
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    order_id int NOT NULL,
    product_id int NOT NULL,
    count int NOT NULL,
    created_on datetime DEFAULT getdate() NOT NULL,
    updated_on datetime,
    CONSTRAINT FK_Order_Line_Items_Orders FOREIGN KEY (order_id) REFERENCES Sales.Orders (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_Order_Line_Items_Products FOREIGN KEY (product_id) REFERENCES Sales.Products (id) ON UPDATE CASCADE
);


CREATE TABLE Sales.Payment_Types
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    name varchar(50),
    description varchar(255)
);

CREATE TABLE Sales.Payments
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    ext_id varchar(64),
    type_id int,
    order_id int,
    amount decimal(18),
    timestamp datetime,
    CONSTRAINT FK_Payments_Payment_Types FOREIGN KEY (type_id) REFERENCES Sales.Payment_Types (id),
    CONSTRAINT FK_Payments_Orders FOREIGN KEY (order_id) REFERENCES Sales.Orders (id)
);
CREATE UNIQUE INDEX Payments_ext_id_uindex ON Sales.Payments (ext_id);

CREATE TABLE School.Courses
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    code varchar(20),
    name varchar(100) NOT NULL,
    description varchar(1000),
    product_id int,
    active bit DEFAULT 0 NOT NULL,
    CONSTRAINT FK_Courses_Products FOREIGN KEY (product_id) REFERENCES Sales.Products (id)
);
CREATE UNIQUE INDEX UQ_Course_Names ON School.Courses (name);
CREATE UNIQUE INDEX UQ_Course_Codes ON School.Courses (code);


CREATE TABLE School.Course_Sessions
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    course_id int NOT NULL,
    session_number int NOT NULL,
    title varchar(50),
    description varchar(500),
    start_datetime datetime,
    CONSTRAINT FK_Course_Sessions_Courses FOREIGN KEY (course_id) REFERENCES School.Courses (id)
);
CREATE UNIQUE INDEX uq_Course_Sessions_titles ON School.Course_Sessions (course_id, title);



  CREATE TABLE School.Course_Registrations
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    course_id int NOT NULL,
    student_id int NOT NULL,
    CONSTRAINT Course_Registrations_Courses_id_fk FOREIGN KEY (course_id) REFERENCES School.Courses (id),
    CONSTRAINT Course_Registrations_Users_id_fk FOREIGN KEY (student_id) REFERENCES People.Users (id)
);


  CREATE TABLE School.Course_Assignments
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    course_session_id int,
    type varchar(50),
    name varchar(100),
    weight int,
    description varchar(1000),
    CONSTRAINT Course_Assignments_Course_Sessions_id_fk FOREIGN KEY (course_session_id) REFERENCES School.Course_Sessions (id) ON UPDATE CASCADE
);

CREATE TABLE School.Course_Assignment_Grades
(
    id uniqueidentifier DEFAULT newid() PRIMARY KEY NOT NULL,
    student_id int NOT NULL,
    assignment_id int NOT NULL,
    attempt int DEFAULT 1 NOT NULL,
    percentage int NOT NULL,
    comments varchar(2000),
    CONSTRAINT Course_Assignment_Grades_Users_id_fk FOREIGN KEY (student_id) REFERENCES People.Users (id) ON UPDATE CASCADE,
    CONSTRAINT Course_Assignment_Grades_Course_Assignments_id_fk FOREIGN KEY (assignment_id) REFERENCES School.Course_Assignments (id) ON UPDATE CASCADE
);
CREATE INDEX Course_Assignment_Grades_id_index ON School.Course_Assignment_Grades (id);
CREATE UNIQUE INDEX Course_Assignment_Grades_student_id_assignment_id_attempt_uindex ON School.Course_Assignment_Grades (student_id, assignment_id, attempt);






CREATE TABLE School.Grade_Letters
(
    letter varchar(2) PRIMARY KEY NOT NULL,
    lowest_percentage int NOT NULL,
    highest_percentage int NOT NULL
);
CREATE UNIQUE INDEX Grade_Letters_highest_percentage_uindex ON School.Grade_Letters (highest_percentage);
CREATE UNIQUE INDEX Grade_Letters_letter_uindex ON School.Grade_Letters (letter);
CREATE UNIQUE INDEX Grade_Letters_lowest_percentage_uindex ON School.Grade_Letters (lowest_percentage);

CREATE TABLE School.Student_Info
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    user_id int NOT NULL,
    major varchar(50),
    start_date date DEFAULT CONVERT([date],getdate()) NOT NULL,
    CONSTRAINT Student_Info_Users_id_fk FOREIGN KEY (user_id) REFERENCES People.Users (id)
);

CREATE TABLE Security.Objects
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    schema_name varchar(50) NOT NULL,
    object_name varchar(100) NOT NULL,
    type varchar(50) NOT NULL
);
CREATE UNIQUE INDEX Objects_object_schema_uindex ON Security.Objects (object_name, schema_name);

CREATE TABLE Security.Permissions
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    name varchar(100) NOT NULL,
    description varchar(500)
);
CREATE UNIQUE INDEX Permissions_name_uindex ON Security.Permissions (name);

CREATE TABLE Security.Roles
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    name varchar(50) NOT NULL,
    description varchar(255)
);
CREATE UNIQUE INDEX Roles_name_uindex ON Security.Roles (name);


CREATE TABLE Security.Object_Role_Permissions
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    object_id int NOT NULL,
    permission_id int NOT NULL,
    role_id int NOT NULL,
    CONSTRAINT Object_Role_Permissions_Objects_id_fk FOREIGN KEY (object_id) REFERENCES Security.Objects (id),
    CONSTRAINT Object_Role_Permissions_Permissions_id_fk FOREIGN KEY (permission_id) REFERENCES Security.Permissions (id),
    CONSTRAINT Object_Role_Permissions_Roles_id_fk FOREIGN KEY (role_id) REFERENCES Security.Roles (id) ON DELETE CASCADE ON UPDATE CASCADE
);




CREATE TABLE Security.User_Roles
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    user_id int,
    role_id int,
    CONSTRAINT User_Roles_Users_id_fk FOREIGN KEY (user_id) REFERENCES People.Users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT User_Roles_Roles_id_fk FOREIGN KEY (role_id) REFERENCES Security.Roles (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE Consulting.Categories
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    name varchar(100) NOT NULL,
    description varchar(1000)
);

CREATE TABLE Consulting.Project_Status
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    status varchar(60) NOT NULL,
    status_display varchar(100)
);
CREATE UNIQUE INDEX Project_Status_status_uindex ON Consulting.Project_Status (status);

CREATE TABLE Consulting.Project_Task_Status
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    status varchar(60) NOT NULL,
    status_display varchar(100)
);
CREATE UNIQUE INDEX Project_Task_Status_status_uindex ON Consulting.Project_Task_Status (status);

CREATE TABLE Consulting.Teams
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    name varchar(100) NOT NULL,
    description varchar(1000),
    specialization varchar(100)
);

CREATE TABLE Consulting.Team_Members
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    team_id int NOT NULL,
    user_id int NOT NULL,
    team_position varchar(255),
    CONSTRAINT Team_Members_Teams_id_fk FOREIGN KEY (team_id) REFERENCES Consulting.Teams (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT Team_Members_Users_id_fk FOREIGN KEY (user_id) REFERENCES People.Users (id) ON DELETE CASCADE ON UPDATE CASCADE
);




CREATE TABLE Accounts.Contacts
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    salutation varchar(20),
    first_name varchar(100) NOT NULL,
    middle_name varchar(100),
    last_name varchar(100) NOT NULL,
    title varchar(255),
    email varchar(100),
    phone_cell varchar(20),
    phone_office varchar(20),
    preferred_method varchar(20),
    notes varchar(1000),
    business_id int
);

CREATE TABLE Accounts.Contact_Log
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    contact_id int,
    contact_method varchar(50) NOT NULL,
    notes varchar(5000),
    timestamp datetime DEFAULT getdate() NOT NULL,
    followup_required bit DEFAULT 0,
    followup_date datetime,
    CONSTRAINT Contact_Log_Contacts_id_fk FOREIGN KEY (contact_id) REFERENCES Accounts.Contacts (id) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE Accounts.Businesses
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    name varchar(100) NOT NULL,
    description varchar(5000)
);
CREATE UNIQUE INDEX Businesses_name_uindex ON Accounts.Businesses (name);


CREATE TABLE Accounts.Business_Locations
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    business_id int NOT NULL,
    name varchar(255) NOT NULL,
    address1 varchar(255),
    address2 varchar(255),
    city varchar(100),
    state varchar(100),
    zip varchar(15),
    country varchar(100),
    phone varchar(20),
    notes varchar(5000),
    CONSTRAINT Business_Locations_Businesses_id_fk FOREIGN KEY (business_id) REFERENCES Accounts.Businesses (id) ON UPDATE CASCADE
);

CREATE TABLE Accounts.Business_Contacts
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    business_id int NOT NULL,
    contact_id int NOT NULL,
    location_id int,
    primary_contact bit DEFAULT 0 NOT NULL,
    notes varchar(500),
    CONSTRAINT Business_Contacts_Businesses_id_fk FOREIGN KEY (business_id) REFERENCES Accounts.Businesses (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT Business_Contacts_Contacts_id_fk FOREIGN KEY (contact_id) REFERENCES Accounts.Contacts (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT Business_Contacts_Business_Locations_id_fk FOREIGN KEY (location_id) REFERENCES Accounts.Business_Locations (id)
);



CREATE TABLE Consulting.Projects
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    business_id int,
    title varchar(255) NOT NULL,
    description varchar(5000),
    category_id int,
    start_date date,
    due_date date,
    primary_contact int,
    team_assigned int,
    status_id int DEFAULT 1 NOT NULL,
    CONSTRAINT Projects_Businesses_id_fk FOREIGN KEY (business_id) REFERENCES Accounts.Businesses (id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT Projects_Categories_id_fk FOREIGN KEY (category_id) REFERENCES Consulting.Categories (id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT Projects_Contacts_id_fk FOREIGN KEY (primary_contact) REFERENCES Accounts.Contacts (id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT Projects_Teams_id_fk FOREIGN KEY (team_assigned) REFERENCES Consulting.Teams (id) ON UPDATE CASCADE,
    CONSTRAINT Projects_Project_Status_id_fk FOREIGN KEY (status_id) REFERENCES Consulting.Project_Status (id) ON UPDATE CASCADE
);


CREATE TABLE Consulting.Project_Tasks
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    project_id int NOT NULL,
    task varchar(255),
    due_date date,
    assigned_user int,
    estimated_hours decimal(18,2),
    notes varchar(1000),
    status_id int DEFAULT 1 NOT NULL,
    CONSTRAINT Project_Tasks_Projects_id_fk FOREIGN KEY (project_id) REFERENCES Consulting.Projects (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT Project_Tasks_Project_Task_Status_id_fk FOREIGN KEY (status_id) REFERENCES Consulting.Project_Task_Status (id) ON UPDATE CASCADE
);

CREATE TABLE Consulting.Project_Task_Time
(
    id int PRIMARY KEY NOT NULL IDENTITY,
    task_id int NOT NULL,
    time decimal(18,2) NOT NULL,
    person_completed int NOT NULL,
    completed_timestamp datetime DEFAULT getdate() NOT NULL,
    notes varchar(1000),
    CONSTRAINT Project_Task_Time_Project_Tasks_id_fk FOREIGN KEY (task_id) REFERENCES Consulting.Project_Tasks (id) ON UPDATE CASCADE,
    CONSTRAINT Project_Task_Time_Users_id_fk FOREIGN KEY (person_completed) REFERENCES People.Users (id) ON UPDATE CASCADE
);





GO

CREATE VIEW School.v_Course_Grades
AS
SELECT        course_id, course_name, student_id, student_username, possible_points, total_points, CAST(total_points / possible_points * 100.0 AS decimal(6, 2)) AS final_percentage,
                             (SELECT        letter
                               FROM            School.Grade_Letters AS l
                               WHERE        (t.total_points / t.possible_points * 100.0 BETWEEN lowest_percentage AND highest_percentage)) AS final_grade
FROM            (SELECT        c.id AS course_id, c.name as course_name, g.student_id, u.username as student_username, SUM(CAST(a.weight * (g.percentage / 100.0) AS decimal(6, 2))) AS total_points, SUM(a.weight) AS possible_points
                          FROM            School.Course_Assignment_Grades AS g INNER JOIN
                      People.Users u on u.id = g.student_id INNER JOIN
                                                    School.Course_Assignments AS a ON a.id = g.assignment_id INNER JOIN
                                                    School.Course_Sessions AS s ON s.id = a.course_session_id INNER JOIN
                                                    School.Courses AS c ON c.id = s.course_id
                          WHERE        (g.attempt =
                                                        (SELECT        MAX(attempt) AS Expr1
                                                          FROM            School.Course_Assignment_Grades AS gr
                                                          WHERE        (student_id = g.student_id) AND (assignment_id = g.assignment_id)))
                          GROUP BY g.student_id, c.id, c.name, u.username) AS t;

GO

CREATE VIEW Consulting.v_Project_Rollup AS
SELECT
    p.id as project_id,
    b.name as customer_name,
    p.title as project_title,
    p.start_date,
    p.due_date,
    ps.status as project_status,
    count(DISTINCT t.id) as total_task_count,
    sum(t.estimated_hours) estimated_hours,

    (
      SELECT count(ptsub.id) FROM Consulting.Project_Tasks ptsub
        INNER JOIN Consulting.Project_Task_Status ptssub on ptssub.id = ptsub.status_id
      WHERE ptsub.project_id = p.id AND ptssub.status = 'completed'
    ) as completed_tasks,

    (
      SELECT SUM(pttsub.time) FROM Consulting.Project_Task_Time pttsub
        INNER JOIN Consulting.Project_Tasks ptsub2 on ptsub2.id = pttsub.task_id
      WHERE ptsub2.project_id = p.id
    ) as actual_hours

  FROM Consulting.Projects p
  INNER JOIN Accounts.Businesses b on b.id = p.business_id
  INNER JOIN Consulting.Project_Tasks t on p.id = t.project_id
  INNER JOIN Consulting.Project_Task_Status pts on t.status_id = pts.id
  INNER JOIN Consulting.Project_Status ps on p.status_id = ps.id

WHERE pts.status NOT IN ('cancelled')

GROUP BY  p.id, p.title, p.start_date, p.due_date, b.name, ps.status;

GO

CREATE FUNCTION [School].[f_GetCourseGrades] (@course_id int)
RETURNS TABLE
AS
RETURN
(
    SELECT
      g.course_id,
      c.name,
      g.student_id,
      u.first_name,
      u.last_name,
      u.username,
      g.final_percentage,
      g.final_grade
    FROM School.v_Course_Grades g

      INNER JOIN School.Courses c on c.id = g.course_id
      INNER JOIN People.Users u on u.id = g.student_id
      INNER JOIN School.Student_Info s ON s.user_id = u.id

    WHERE course_id = @course_id
);

GO


EXEC sp_msforeachtable 'ALTER TABLE ? NOCHECK CONSTRAINT all';


INSERT INTO dbo.access_keys (kid, x5c, timestamp) VALUES ('1LTMzakihiRla_8z2BEJVXeWMqo', 'MIIDYDCCAkigAwIBAgIJAIB4jVVJ3BeuMA0GCSqGSIb3DQEBCwUAMCkxJzAlBgNVBAMTHkxpdmUgSUQgU1RTIFNpZ25pbmcgUHVibGljIEtleTAeFw0xNjA0MDUxNDQzMzVaFw0yMTA0MDQxNDQzMzVaMCkxJzAlBgNVBAMTHkxpdmUgSUQgU1RTIFNpZ25pbmcgUHVibGljIEtleTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAN7CnCUg+HB8E2OY2JuZTRM6pN8DjdQmgETuRKFiELebvl/mZ1utQokWaVqX4SsEWFNed7XeUfKfgCQ1k/jXBCdDB1NyoriBD5m4CaUu6pMlEbh0bnAzNpTuC+F3yJb4s0K5ciZyXHdPbixK2RDvT3Lyogx6KIrNgeHxo+wrAHuVGYBMSbwAg5bEEpLtiVMg5mVGT9v7v4ArqjSTOaNIVuxSkSNzBaGArKkhJh557pridoxdERw0VkKaHYZQEerii6Ilh6hzt2hu7HMWDI+XDM5KXtSqvw5ucsh5DAkifyWkYcMphX1f0lpKnd5/llWU9AAyxKtETYvbameN56FJzv8CAwEAAaOBijCBhzAdBgNVHQ4EFgQU9IdLLpbC2S8Wn1MCXsdtFac9SRYwWQYDVR0jBFIwUIAU9IdLLpbC2S8Wn1MCXsdtFac9SRahLaQrMCkxJzAlBgNVBAMTHkxpdmUgSUQgU1RTIFNpZ25pbmcgUHVibGljIEtleYIJAIB4jVVJ3BeuMAsGA1UdDwQEAwIBxjANBgkqhkiG9w0BAQsFAAOCAQEAXk0sQAib0PGqvwELTlflQEKS++vqpWYPW/2gCVCn5shbyP1J7z1nT8kE/ZDVdl3LvGgTMfdDHaRF5ie5NjkTHmVOKbbHaWpTwUFbYAFBJGnx+s/9XSdmNmW9GlUjdpd6lCZxsI6888r0ptBgKINRRrkwMlq3jD1U0kv4JlsIhafUIOqGi4+hIDXBlY0F/HJPfUU75N885/r4CCxKhmfh3PBM35XOch/NGC67fLjqLN+TIWLoxnvil9m3jRjqOA9u50JUeDGZABIYIMcAdLpI2lcfru4wXcYXuQul22nAR7yOyGKNOKULoOTE4t4AeGRqCogXSxZgaTgKSBhvhE+MGg==', '2018-10-29 16:10:30.530');
INSERT INTO dbo.access_keys (kid, x5c, timestamp) VALUES ('2S4SCVGs8Sg9LS6AqLIq6DpW-g8', 'MIIDKDCCAhCgAwIBAgIQBHJvVNxP1oZO4HYKh+rypDANBgkqhkiG9w0BAQsFADAjMSEwHwYDVQQDExhsb2dpbi5taWNyb3NvZnRvbmxpbmUudXMwHhcNMTYxMTE2MDgwMDAwWhcNMTgxMTE2MDgwMDAwWjAjMSEwHwYDVQQDExhsb2dpbi5taWNyb3NvZnRvbmxpbmUudXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQChn5BCs24Hh6L0BNitPrV5s+2/DBhaeytOmnghJKnqeJlhv3ZczShRM2Cp38LW8Y3wn7L3AJtolaSkF/joKN1l6GupzM+HOEdq7xZxFehxIHW7+25mG/WigBnhsBzLv1SR4uIbrQeS5M0kkLwJ9pOnVH3uzMGG6TRXPnK3ivlKl97AiUEKdlRjCQNLXvYf1ZqlC77c/ZCOHSX4kvIKR2uG+LNlSTRq2rn8AgMpFT4DSlEZz4RmFQvQupQzPpzozaz/gadBpJy/jgDmJlQMPXkHp7wClvbIBGiGRaY6eZFxNV96zwSR/GPNkTObdw2S8/SiAgvIhIcqWTPLY6aVTqJfAgMBAAGjWDBWMFQGA1UdAQRNMEuAEDUj0BrjP0RTbmoRPTRMY3WhJTAjMSEwHwYDVQQDExhsb2dpbi5taWNyb3NvZnRvbmxpbmUudXOCEARyb1TcT9aGTuB2Cofq8qQwDQYJKoZIhvcNAQELBQADggEBAGnLhDHVz2gLDiu9L34V3ro/6xZDiSWhGyHcGqky7UlzQH3pT5so8iF5P0WzYqVtogPsyC2LPJYSTt2vmQugD4xlu/wbvMFLcV0hmNoTKCF1QTVtEQiAiy0Aq+eoF7Al5fV1S3Sune0uQHimuUFHCmUuF190MLcHcdWnPAmzIc8fv7quRUUsExXmxSX2ktUYQXzqFyIOSnDCuWFm6tpfK5JXS8fW5bpqTlrysXXz/OW/8NFGq/alfjrya4ojrOYLpunGriEtNPwK7hxj1AlCYEWaRHRXaUIW1ByoSff/6Y6+ZhXPUe0cDlNRt/qIz5aflwO7+W8baTS4O8m/icu7ItE=', '2018-10-29 16:10:30.530');
INSERT INTO dbo.access_keys (kid, x5c, timestamp) VALUES ('i6lGk3FZzxRcUb2C3nEQ7syHJlY', 'MIIDBTCCAe2gAwIBAgIQWOnQG5bZiIFAbuSzrxjvbDANBgkqhkiG9w0BAQsFADAtMSswKQYDVQQDEyJhY2NvdW50cy5hY2Nlc3Njb250cm9sLndpbmRvd3MubmV0MB4XDTE4MDgwMTAwMDAwMFoXDTIwMDgwMTAwMDAwMFowLTErMCkGA1UEAxMiYWNjb3VudHMuYWNjZXNzY29udHJvbC53aW5kb3dzLm5ldDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALM5fWyNaHwYqnCfxVBEqhfDHKWnNhqmVHH9McYM+bsRs9xl7DkuwsbvD6vOb4rANtVOuly+Eiv/fYhhwn8UFE4kzTRXWREIH+uoFCUDyItwtT8vNn253XzXuFiIeHEkSoHM4vckdP+G7lPEQRZtQjV2TTB/76eBrtnLdP/AAgp4KE4/kjSuKEvn4lwiKtQWYPmCBPh18erwmIFwBramK5CzNKT31ycTYpi3LZpjbIljzW7gzvX9P6/U+dmTmpOufKzJbAnnXKveMcRUZl7wXWxQpyKGAZ0q+8FKIcVbrO0mlbyPuQGynAOofsnhqPQVeO4lbsEgcmL9QXWdxdn6kYsCAwEAAaMhMB8wHQYDVR0OBBYEFMXZ+lF0trQrd4uWs+lQ7h0mls63MA0GCSqGSIb3DQEBCwUAA4IBAQAA/TRp5Cx6WWA04dCvoJNTd80KSO8uGrM2V4VhTlIU1zf8cjMocmBXSA0v+P3onkBHHwnMJs6fWuwcBL4vhqZ4jonPcgl3tUIGgDwywM3V4mC1JLZJXlBZASsAuKmHm6qwge6RVs+Ub1fcpOIViMgW/1Wj1tZm/v/NGHQ+EEJV1UXtE5OnFnzD+9TrfOTI/l855ryKNS6I+XNbtcIJUL87OWkMRgNgjOhBnmldA27sWWkFWfvxB0J7FvOkyWaLlyNe/jqCL3jYXzz0XBNGMbatmeKS3fgIihVUJ32Vt7GXI4ed0HuvhhZ6d+fvMnBPWdfteu018YHKeXyk/c2aegcj', '2018-10-29 16:10:30.530');
INSERT INTO dbo.access_keys (kid, x5c, timestamp) VALUES ('M6pX7RHoraLsprfJeRCjSxuURhc', 'MIIC8TCCAdmgAwIBAgIQfEWlTVc1uINEc9RBi6qHMjANBgkqhkiG9w0BAQsFADAjMSEwHwYDVQQDExhsb2dpbi5taWNyb3NvZnRvbmxpbmUudXMwHhcNMTgxMDE0MDAwMDAwWhcNMjAxMDE0MDAwMDAwWjAjMSEwHwYDVQQDExhsb2dpbi5taWNyb3NvZnRvbmxpbmUudXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDEdJxkw+jwWJ+gNyuCdxZDuYYm2IqGuyGjT64U9eD452dEpi51MPv4GrUwONypZF/ch8NWMUmLFBLrUzvCb3AAsOP76Uu4kn2MNQBZMDfFa9AtuIEz6CpTSyPiZzaVabqc9+qXJh5a1BxILSxiQuVrI2BfiegoNpeK+FU6ntvAZervsxHN4vj72qtFgqO7Md9kvuQ0EKyo7Xzk/Q0jYm2bD4SypiysJoex81EZGtO9QdSreFZrzn2Qr/m413tN5jZkTApUPx7MKZJ9Hn1nPLFO24+mQJIdL061S9LeapNiK3vepy+muOXdHyGmNctvyh+1+laveEVF2nGvC6hAQ7hBAgMBAAGjITAfMB0GA1UdDgQWBBQ5TKadw06O0cvXrQbXW0Nb3M3h/DANBgkqhkiG9w0BAQsFAAOCAQEAI48JaFtwOFcYS/3pfS5+7cINrafXAKTL+/+he4q+RMx4TCu/L1dl9zS5W1BeJNO2GUznfI+b5KndrxdlB6qJIDf6TRHh6EqfA18oJP5NOiKhU4pgkF2UMUw4kjxaZ5fQrSoD9omjfHAFNjradnHA7GOAoF4iotvXDWDBWx9K4XNZHWvD11Td66zTg5IaEQDIZ+f8WS6nn/98nAVMDtR9zW7Te5h9kGJGfe6WiHVaGRPpBvqC4iypGHjbRwANwofZvmp5wP08hY1CsnKY5tfP+E2k/iAQgKKa6QoxXToYvP7rsSkglak8N5g/+FJGnq4wP6cOzgZpjdPMwaVt5432GA==', '2018-10-29 16:10:30.530');
INSERT INTO dbo.access_keys (kid, x5c, timestamp) VALUES ('wULmYfsqdQuWtV_-hxVtDJJZM4Q', 'MIIDBTCCAe2gAwIBAgIQKGZsKfAUzaJHVantyrwVdzANBgkqhkiG9w0BAQsFADAtMSswKQYDVQQDEyJhY2NvdW50cy5hY2Nlc3Njb250cm9sLndpbmRvd3MubmV0MB4XDTE4MTAwMTAwMDAwMFoXDTIwMTAwMTAwMDAwMFowLTErMCkGA1UEAxMiYWNjb3VudHMuYWNjZXNzY29udHJvbC53aW5kb3dzLm5ldDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAJrUafn5X9QkUDRCyYHsbK3tGgcQ1GiHEcrJmmAcVsYB9K9nH3cKHu3nVraOazG5EpF1ccVfPevQs+lg5l90pjkgHVTLt51V3CDlSG5yCOUc/Vuo+WKyJ8j/xT0P17meFjH17iLRlzfG6H3pdXoQW4UvWBL4BaXzwQr0g9SVQq1A/2GbVxRj0tPvvZ3E8TTxHLKNzSa6fhQxEmfH/J340QdZn3Hj8BnSrLs53S9xRj2y5/IaTuGV5/ysrCJmuFwZ6r27NF9AKlhBv4O2my+txzPQou0Duz2OkZuVw3IviofnKJ1nCl7+c7eTTioIzonqa+VTd240pfHKU9oLWUrRbUsCAwEAAaMhMB8wHQYDVR0OBBYEFFve38eLO2PUMXcqbBC/YDaayIbrMA0GCSqGSIb3DQEBCwUAA4IBAQCWhrem9A2NrOiGesdzbDy2K3k0oxjWMlM/ZGIfPMPtZl4uIzhc+vdDVVFSeV8SKTOEIMjMOJTQ3GJpZEHYlyM7UsGWiMSXqzG5HUxbkPvuEFHx7cl9Ull3AEymB2oVPC9DPtLUXPyDH898QgEEVhAEI+JZc1Yd6mAlY/5nOw5m2Yqm+84JOPWLgFDqfVmz/MH27LS1rnzzc+0hhcm/Nv/x7FmpOeRfh00BjCA4PogJlpjl/z/6+GTYcYFsvKE3jmmXka8tQbBOHgAlMnamFA8xGeDok6QaxOELu8NSWzvyZXM2lJK5WFQPHF2hjnNXs6+RxOovG55Ybpo52c2frhNZ', '2018-10-29 16:10:30.530');
INSERT INTO dbo.access_keys (kid, x5c, timestamp) VALUES ('xP_zn6I1YkXcUUmlBoPuXTGsaxk', 'MIIDYDCCAkigAwIBAgIJAJzCyTLC+DjJMA0GCSqGSIb3DQEBCwUAMCkxJzAlBgNVBAMTHkxpdmUgSUQgU1RTIFNpZ25pbmcgUHVibGljIEtleTAeFw0xNjA3MTMyMDMyMTFaFw0yMTA3MTIyMDMyMTFaMCkxJzAlBgNVBAMTHkxpdmUgSUQgU1RTIFNpZ25pbmcgUHVibGljIEtleTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANqVmrWn3m95gdAO9/mflEcK7m8A5XVr4r0btfRjQutTATm932eCOJcoVITbku2iUTySBQes5fhxcdZx6RoO+AR3tM/bJDgx2GRAAndoigSHreeQ3XbiHtQl7GfZHshJjWhr910uDlv+yAzNI+cEfTh7oFIn5B9z0CHBb02Gpx5l6I05oIyDlqM+8tKJKcjffhW2Uc3X8fno/p42o+YEfK+vZrsDucJU2yBttFeNEes2IfkDgzFKWssu4wxhRsOb4xR/QmtLS2XDUKv/yShP5qftHzgzgEKGzqwVKAEhZMgnfeQVVuyia53b1u6H0Vmgk7mWTsuUp1OGrG4zK4MRWRcCAwEAAaOBijCBhzAdBgNVHQ4EFgQU11z579/IePwuc4WBdN4L0ljG4CUwWQYDVR0jBFIwUIAU11z579/IePwuc4WBdN4L0ljG4CWhLaQrMCkxJzAlBgNVBAMTHkxpdmUgSUQgU1RTIFNpZ25pbmcgUHVibGljIEtleYIJAJzCyTLC+DjJMAsGA1UdDwQEAwIBxjANBgkqhkiG9w0BAQsFAAOCAQEAiASLEpQseGNahE+9f9PQgmX3VgjJerNjXr1zXWXDJfFE31DxgsxddjcIgoBL9lwegOHHvwpzK1ecgH45xcJ0Z/40OgY8NITqXbQRfdgLrEGJCoyOQEbjb5PW5k2aOdn7LBxvDsH6Y8ax26v+EFMPh3G+xheh6bfoIRSK1b+44PfoDZoJ9NfJibOZ4Cq+wt/yOvpMYQDB/9CNo18wmA3RCLYjf2nAc7RO0PDYHSIq5QDWV+1awmXDKgIdRpYPpRtn9KFXQkpCeEc/lDTG+o6n7nC40wyjioyR6QmHGvNkMR4VfSoTKCTnFATyDpI1bqU2K7KNjUEsCYfwybFB8d6mjQ==', '2018-10-29 16:10:30.530');


INSERT INTO dbo.casbin_rule (p_type, v0, v1, v2, v3, v4, v5) VALUES ('p', 'Employees', '/users/:user_id', 'GET', '', '', '');
INSERT INTO dbo.casbin_rule (p_type, v0, v1, v2, v3, v4, v5) VALUES ('p', 'Employees', '/courses/:course_id/sessions/:session_id', 'GET', '', '', '');
INSERT INTO dbo.casbin_rule (p_type, v0, v1, v2, v3, v4, v5) VALUES ('p', 'Employees', '/courses', 'GET', '', '', '');
INSERT INTO dbo.casbin_rule (p_type, v0, v1, v2, v3, v4, v5) VALUES ('p', 'Public', '/publicendpoint', 'GET', '', '', '');
INSERT INTO dbo.casbin_rule (p_type, v0, v1, v2, v3, v4, v5) VALUES ('g', 'Brandon.Taylor@fortisureit.com', 'Employees', '', '', '', '');
INSERT INTO dbo.casbin_rule (p_type, v0, v1, v2, v3, v4, v5) VALUES ('g', 'Test.Patterson@fortisureit.com', 'Public', '', '', '', '');
INSERT INTO dbo.casbin_rule (p_type, v0, v1, v2, v3, v4, v5) VALUES ('p', 'Employees', '/protected', 'GET', null, null, null);
INSERT INTO dbo.casbin_rule (p_type, v0, v1, v2, v3, v4, v5) VALUES ('g', 'Brad.Magyar@fortisureit.com', 'Admins', null, null, null, null);
INSERT INTO dbo.casbin_rule (p_type, v0, v1, v2, v3, v4, v5) VALUES ('p', 'Admins', '/*', 'GET|PUT|POST|DELETE', null, null, null);


SET IDENTITY_INSERT People.Users ON;
INSERT INTO People.Users (id, first_name, last_name, username, ad_object_id) VALUES (1, 'Brad', 'Magyar', 'Brad.Magyar', 'f1936400-aebc-4468-b4f7-401dd10eb06a');
INSERT INTO People.Users (id, first_name, last_name, username, ad_object_id) VALUES (2, 'Scott', 'Arnold', 'Scott.Arnold', '3e23f737-f866-4e45-90b8-dcef8b71954e');
INSERT INTO People.Users (id, first_name, last_name, username, ad_object_id) VALUES (3, 'Brandon', 'Taylor', 'Brandon.Taylor', 'dd159bf5-d9ea-458f-8009-dff30c89f4f3');
INSERT INTO People.Users (id, first_name, last_name, username, ad_object_id) VALUES (4, 'Test', 'Patterson', 'Test.Patterson', null);
INSERT INTO People.Users (id, first_name, last_name, username, ad_object_id) VALUES (5, 'Jonny', 'Karate', 'Jonny.Karate', null);
INSERT INTO People.Users (id, first_name, last_name, username, ad_object_id) VALUES (6, 'Bill', 'Brasky', 'Bill.Brasky', null);
SET IDENTITY_INSERT People.Users OFF;


SET IDENTITY_INSERT People.User_Info ON;
INSERT INTO People.User_Info (id, user_id, preferred_name, gender, date_of_birth, email, phone_primary, phone_secondary, address1, address2, city, state, zip, bio) VALUES (1, 1, 'Brad', 'Male', '1989-03-14', 'brad.magyar@fortisureit.com', '3306344281', null, '346 Center Rd', null, 'New Franklin', 'OH', '44319', 'Just a simple man trying to make his way in the universe');
INSERT INTO People.User_Info (id, user_id, preferred_name, gender, date_of_birth, email, phone_primary, phone_secondary, address1, address2, city, state, zip, bio) VALUES (2, 2, 'Scott', 'Male', null, 'scott.arnold@fortisureit.com', null, null, null, null, null, null, null, null);
INSERT INTO People.User_Info (id, user_id, preferred_name, gender, date_of_birth, email, phone_primary, phone_secondary, address1, address2, city, state, zip, bio) VALUES (3, 3, 'Brandon', 'Male', null, 'brandon.taylor@fortisureit.com', null, null, null, null, null, null, null, null);
SET IDENTITY_INSERT People.User_Info OFF;


SET IDENTITY_INSERT Sales.Orders ON;
INSERT INTO Sales.Orders (id, user_id, status, notes, created_on, updated_on) VALUES (1, 4, 'open', 'Course registrations', '2018-08-04 05:23:46.527', null);
INSERT INTO Sales.Orders (id, user_id, status, notes, created_on, updated_on) VALUES (2, 5, 'paid', 'Single course registration', '2018-08-04 05:38:13.367', null);
SET IDENTITY_INSERT Sales.Orders OFF;



SET IDENTITY_INSERT Sales.Order_Line_Items ON;
INSERT INTO Sales.Order_Line_Items (id, order_id, product_id, count, created_on, updated_on) VALUES (1, 1, 1, 1, '2018-08-04 05:24:09.107', null);
INSERT INTO Sales.Order_Line_Items (id, order_id, product_id, count, created_on, updated_on) VALUES (2, 1, 2, 4, '2018-08-04 05:29:31.030', null);
INSERT INTO Sales.Order_Line_Items (id, order_id, product_id, count, created_on, updated_on) VALUES (3, 2, 2, 2, '2018-08-04 05:38:55.953', null);
SET IDENTITY_INSERT Sales.Order_Line_Items OFF;


SET IDENTITY_INSERT Sales.Payment_Types ON;
INSERT INTO Sales.Payment_Types (id, name, description) VALUES (1, 'Credit Card', 'Credit card payments');
INSERT INTO Sales.Payment_Types (id, name, description) VALUES (2, 'Cash', 'Cash payments');
INSERT INTO Sales.Payment_Types (id, name, description) VALUES (3, 'Check', 'Check payments');
SET IDENTITY_INSERT Sales.Payment_Types OFF;

SET IDENTITY_INSERT Sales.Products ON;
INSERT INTO Sales.Products (id, price, name, category, description) VALUES (1, 50.00, 'IT101 Course', 'Courses', 'IT101 - IT Essentials training course');
INSERT INTO Sales.Products (id, price, name, category, description) VALUES (2, 99.99, 'IT201', 'Courses', 'IT201 - Advanced IT training course');
SET IDENTITY_INSERT Sales.Products OFF;

INSERT INTO School.Course_Assignment_Grades (id, student_id, assignment_id, attempt, percentage, comments) VALUES ('E9B30485-3795-431D-A6FF-066B73975699', 4, 4, 1, 93, '');
INSERT INTO School.Course_Assignment_Grades (id, student_id, assignment_id, attempt, percentage, comments) VALUES ('7BB3BA14-7BDE-4AFA-9C22-099CF05E3782', 4, 3, 1, 91, null);
INSERT INTO School.Course_Assignment_Grades (id, student_id, assignment_id, attempt, percentage, comments) VALUES ('D6AE8979-47ED-4DF5-B942-280B98429508', 5, 1, 3, 43, 'WHY DID YOU TRY AGAIN?!');
INSERT INTO School.Course_Assignment_Grades (id, student_id, assignment_id, attempt, percentage, comments) VALUES ('C490A5C1-6165-49C3-8693-3189A945CEF8', 4, 5, 1, 86, null);
INSERT INTO School.Course_Assignment_Grades (id, student_id, assignment_id, attempt, percentage, comments) VALUES ('E3E5E329-6741-465B-BEB5-51807EC21C6A', 4, 4, 2, 100, null);
INSERT INTO School.Course_Assignment_Grades (id, student_id, assignment_id, attempt, percentage, comments) VALUES ('A952C7DA-E942-459C-AC83-565FB16DC7CA', 4, 2, 1, 84, null);
INSERT INTO School.Course_Assignment_Grades (id, student_id, assignment_id, attempt, percentage, comments) VALUES ('19A15453-1D13-4B3B-A596-5E8C6D91BACA', 4, 4, 3, 90, null);
INSERT INTO School.Course_Assignment_Grades (id, student_id, assignment_id, attempt, percentage, comments) VALUES ('EE44F529-4042-42D2-A589-628B7E6FBDEF', 5, 3, 1, 79, null);
INSERT INTO School.Course_Assignment_Grades (id, student_id, assignment_id, attempt, percentage, comments) VALUES ('86DDA3A5-7384-49EC-A1FE-936E4EBB9BAA', 5, 2, 1, 88, null);
INSERT INTO School.Course_Assignment_Grades (id, student_id, assignment_id, attempt, percentage, comments) VALUES ('E450D213-4B1B-43F9-A300-97942AF21612', 5, 4, 1, 98, null);
INSERT INTO School.Course_Assignment_Grades (id, student_id, assignment_id, attempt, percentage, comments) VALUES ('E0909C27-C459-4C8B-A67E-BE0A15372BA6', 5, 5, 1, 68, null);
INSERT INTO School.Course_Assignment_Grades (id, student_id, assignment_id, attempt, percentage, comments) VALUES ('520AB744-007E-4001-8F74-DE2C6A79591F', 5, 1, 1, 50, 'Not so good');
INSERT INTO School.Course_Assignment_Grades (id, student_id, assignment_id, attempt, percentage, comments) VALUES ('FEF8DC0B-E16A-4D3F-A349-ED92031D75B9', 5, 1, 2, 88, 'Much better');
INSERT INTO School.Course_Assignment_Grades (id, student_id, assignment_id, attempt, percentage, comments) VALUES ('40ACB844-45FC-4B1A-B450-FC7B9F996F27', 4, 1, 1, 95, 'Well done');

SET IDENTITY_INSERT School.Course_Assignments ON;
INSERT INTO School.Course_Assignments (id, course_session_id, type, name, weight, description) VALUES (1, 1, 'Quiz', 'Unit 1 Quiz', 10, 'Quiz of basic computer concepts');
INSERT INTO School.Course_Assignments (id, course_session_id, type, name, weight, description) VALUES (2, 2, 'Lab', 'Unit 2 Lab', 20, 'Unit 2 hands on lab');
INSERT INTO School.Course_Assignments (id, course_session_id, type, name, weight, description) VALUES (3, 3, 'Lab', 'Unit 3 Lab', 15, 'Unit 3 lab');
INSERT INTO School.Course_Assignments (id, course_session_id, type, name, weight, description) VALUES (4, 4, 'Test', 'Unit 1 Test', 40, null);
INSERT INTO School.Course_Assignments (id, course_session_id, type, name, weight, description) VALUES (5, 5, 'Lab', 'Unit 2 Lab', 25, null);
SET IDENTITY_INSERT School.Course_Assignments OFF;

SET IDENTITY_INSERT School.Course_Registrations ON;
INSERT INTO School.Course_Registrations (id, course_id, student_id) VALUES (1, 1, 4);
INSERT INTO School.Course_Registrations (id, course_id, student_id) VALUES (2, 1, 5);
INSERT INTO School.Course_Registrations (id, course_id, student_id) VALUES (3, 2, 5);
INSERT INTO School.Course_Registrations (id, course_id, student_id) VALUES (4, 2, 4);
INSERT INTO School.Course_Registrations (id, course_id, student_id) VALUES (5, 3, 5);
SET IDENTITY_INSERT School.Course_Registrations OFF ;

SET IDENTITY_INSERT School.Course_Sessions ON;
INSERT INTO School.Course_Sessions (id, course_id, session_number, title, description, start_datetime) VALUES (1, 1, 1, 'Introduction', 'Basic IT concepts', '2018-09-13 18:12:50.963');
INSERT INTO School.Course_Sessions (id, course_id, session_number, title, description, start_datetime) VALUES (2, 1, 2, 'Systems Architecture', 'Computer systems introduction', '2018-09-13 18:12:50.963');
INSERT INTO School.Course_Sessions (id, course_id, session_number, title, description, start_datetime) VALUES (3, 1, 3, 'Hardware Overview', 'Overview of computer hardware', '2018-09-13 18:12:50.963');
INSERT INTO School.Course_Sessions (id, course_id, session_number, title, description, start_datetime) VALUES (4, 2, 1, 'Test', 'Testing 123', '2018-09-13 18:12:50.963');
INSERT INTO School.Course_Sessions (id, course_id, session_number, title, description, start_datetime) VALUES (5, 3, 1, 'Databases', 'Overview of database types and concepts', '2018-09-13 18:12:50.963');
INSERT INTO School.Course_Sessions (id, course_id, session_number, title, description, start_datetime) VALUES (6, 3, 2, 'SQL Intro', 'Introduction to SQL language - SQL Server focus', '2018-09-13 18:12:50.963');
SET IDENTITY_INSERT School.Course_Sessions OFF;


SET IDENTITY_INSERT School.Courses ON;
INSERT INTO School.Courses (id, code, name, description, product_id, active) VALUES (1, 'IT101', 'IT Essentials', 'A survey of essential IT concepts', null, 1);
INSERT INTO School.Courses (id, code, name, description, product_id, active) VALUES (2, 'IT201', 'Advanced IT', 'A deeper dive into essential IT concepts', null, 1);
INSERT INTO School.Courses (id, code, name, description, product_id, active) VALUES (3, 'DB101', 'Database Fundamentals', 'Database administration essentials', null, 1);
INSERT INTO School.Courses (id, code, name, description, product_id, active) VALUES (4, 'PR101', 'Programming Fundamentals', 'Essential programming concepts', null, 1);
SET IDENTITY_INSERT School.Courses OFF;

INSERT INTO School.Grade_Letters (letter, lowest_percentage, highest_percentage) VALUES ('A', 94, 99);
INSERT INTO School.Grade_Letters (letter, lowest_percentage, highest_percentage) VALUES ('A+', 100, 200);
INSERT INTO School.Grade_Letters (letter, lowest_percentage, highest_percentage) VALUES ('A-', 90, 94);
INSERT INTO School.Grade_Letters (letter, lowest_percentage, highest_percentage) VALUES ('B', 83, 85);
INSERT INTO School.Grade_Letters (letter, lowest_percentage, highest_percentage) VALUES ('B+', 86, 89);
INSERT INTO School.Grade_Letters (letter, lowest_percentage, highest_percentage) VALUES ('B-', 80, 82);
INSERT INTO School.Grade_Letters (letter, lowest_percentage, highest_percentage) VALUES ('C', 73, 75);
INSERT INTO School.Grade_Letters (letter, lowest_percentage, highest_percentage) VALUES ('C+', 76, 79);
INSERT INTO School.Grade_Letters (letter, lowest_percentage, highest_percentage) VALUES ('C-', 70, 72);
INSERT INTO School.Grade_Letters (letter, lowest_percentage, highest_percentage) VALUES ('D', 63, 65);
INSERT INTO School.Grade_Letters (letter, lowest_percentage, highest_percentage) VALUES ('D+', 66, 69);
INSERT INTO School.Grade_Letters (letter, lowest_percentage, highest_percentage) VALUES ('D-', 60, 62);
INSERT INTO School.Grade_Letters (letter, lowest_percentage, highest_percentage) VALUES ('F', 0, 59);

SET IDENTITY_INSERT School.Student_Info ON;
INSERT INTO School.Student_Info (id, user_id, major, start_date) VALUES (1, 4, 'Database', '2018-07-14');
INSERT INTO School.Student_Info (id, user_id, major, start_date) VALUES (2, 5, 'Programming', '2018-07-14');
SET IDENTITY_INSERT School.Student_Info OFF;

SET IDENTITY_INSERT Security.Objects ON;
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (1, 'Company', 'Departments', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (2, 'Company', 'Employee_Info', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (3, 'Company', 'Positions', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (4, 'People', 'User_Info', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (5, 'Security', 'User_Roles', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (6, 'People', 'Users', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (7, 'Sales', 'Order_Line_Items', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (8, 'Sales', 'Orders', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (9, 'Sales', 'Payment_Types', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (10, 'Sales', 'Payments', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (11, 'Sales', 'Products', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (12, 'School', 'Course_Assignment_Grades', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (13, 'School', 'Course_Assignments', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (14, 'School', 'Course_Registrations', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (15, 'School', 'Course_Sessions', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (16, 'School', 'Courses', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (17, 'School', 'Grade_Letters', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (18, 'School', 'Student_Info', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (19, 'Security', 'Objects', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (20, 'Security', 'Permissions', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (21, 'Security', 'Roles', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (22, 'School', 'v_Course_Grades', 'View');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (23, 'Security', 'Object_Role_Permissions', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (24, 'Accounts', 'Business_Contacts', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (25, 'Accounts', 'Business_Locations', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (26, 'Accounts', 'Businesses', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (27, 'Accounts', 'Contact_Log', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (28, 'Accounts', 'Contacts', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (29, 'Consulting', 'Categories', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (30, 'Consulting', 'Project_Status', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (31, 'Consulting', 'Project_Task_Status', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (32, 'Consulting', 'Project_Task_Time', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (33, 'Consulting', 'Project_Tasks', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (34, 'Consulting', 'Projects', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (35, 'Consulting', 'Team_Members', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (36, 'Consulting', 'Teams', 'Table');
INSERT INTO Security.Objects (id, schema_name, object_name, type) VALUES (37, 'Consulting', 'v_Project_Rollup', 'View');
SET IDENTITY_INSERT Security.Objects OFF;

SET IDENTITY_INSERT Security.Object_Role_Permissions ON;
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (1, 1, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (2, 1, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (3, 1, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (4, 1, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (5, 2, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (6, 2, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (7, 2, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (8, 2, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (9, 3, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (10, 3, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (11, 3, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (12, 3, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (13, 4, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (14, 4, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (15, 4, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (16, 4, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (17, 5, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (18, 5, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (19, 5, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (20, 5, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (21, 6, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (22, 6, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (23, 6, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (24, 6, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (25, 7, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (26, 7, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (27, 7, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (28, 7, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (29, 8, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (30, 8, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (31, 8, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (32, 8, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (33, 9, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (34, 9, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (35, 9, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (36, 9, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (37, 10, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (38, 10, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (39, 10, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (40, 10, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (41, 11, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (42, 11, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (43, 11, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (44, 11, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (45, 12, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (46, 12, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (47, 12, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (48, 12, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (49, 13, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (50, 13, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (51, 13, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (52, 13, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (53, 14, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (54, 14, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (55, 14, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (56, 14, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (57, 15, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (58, 15, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (59, 15, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (60, 15, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (61, 16, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (62, 16, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (63, 16, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (64, 16, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (65, 17, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (66, 17, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (67, 17, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (68, 17, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (69, 18, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (70, 18, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (71, 18, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (72, 18, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (73, 19, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (74, 19, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (75, 19, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (76, 19, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (77, 20, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (78, 20, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (79, 20, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (80, 20, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (81, 21, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (82, 21, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (83, 21, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (84, 21, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (85, 22, 1, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (86, 22, 2, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (87, 22, 3, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (88, 22, 4, 1);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (89, 1, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (90, 2, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (91, 3, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (92, 4, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (93, 5, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (94, 6, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (95, 7, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (96, 8, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (97, 9, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (98, 10, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (99, 11, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (100, 12, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (101, 13, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (102, 14, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (103, 15, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (104, 16, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (105, 17, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (106, 18, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (107, 19, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (108, 20, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (109, 21, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (110, 22, 1, 2);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (111, 12, 1, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (112, 13, 1, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (113, 14, 1, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (114, 15, 1, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (115, 16, 1, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (116, 17, 1, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (117, 18, 1, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (118, 22, 1, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (119, 12, 2, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (120, 13, 2, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (121, 14, 2, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (122, 15, 2, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (123, 16, 2, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (124, 17, 2, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (125, 18, 2, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (126, 22, 2, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (127, 12, 3, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (128, 13, 3, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (129, 14, 3, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (130, 15, 3, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (131, 16, 3, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (132, 17, 3, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (133, 18, 3, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (134, 22, 3, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (135, 12, 4, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (136, 13, 4, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (137, 14, 4, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (138, 15, 4, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (139, 16, 4, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (140, 17, 4, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (141, 18, 4, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (142, 22, 4, 3);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (143, 12, 1, 4);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (144, 13, 1, 4);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (145, 14, 1, 4);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (146, 15, 1, 4);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (147, 16, 1, 4);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (148, 17, 1, 4);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (149, 18, 1, 4);
INSERT INTO Security.Object_Role_Permissions (id, object_id, permission_id, role_id) VALUES (150, 22, 1, 4);
SET IDENTITY_INSERT Security.Object_Role_Permissions OFF;

SET IDENTITY_INSERT Security.Permissions ON;
INSERT INTO Security.Permissions (id, name, description) VALUES (1, 'GET', 'Read entries');
INSERT INTO Security.Permissions (id, name, description) VALUES (2, 'PUT', 'Create new entries');
INSERT INTO Security.Permissions (id, name, description) VALUES (3, 'POST', 'Update existing entries');
INSERT INTO Security.Permissions (id, name, description) VALUES (4, 'DELETE', 'Delete existing entries');
SET IDENTITY_INSERT Security.Permissions OFF;

SET IDENTITY_INSERT Security.Roles ON;
INSERT INTO Security.Roles (id, name, description) VALUES (1, 'Admin', 'Global Administrator');
INSERT INTO Security.Roles (id, name, description) VALUES (2, 'Employee', 'Standard Employee');
INSERT INTO Security.Roles (id, name, description) VALUES (3, 'Instructor', 'Course Instructor');
INSERT INTO Security.Roles (id, name, description) VALUES (4, 'Student', 'Course Student');
SET IDENTITY_INSERT Security.Roles OFF;

SET IDENTITY_INSERT Security.User_Roles ON;
INSERT INTO Security.User_Roles (id, user_id, role_id) VALUES (1, 1, 1);
INSERT INTO Security.User_Roles (id, user_id, role_id) VALUES (2, 1, 3);
INSERT INTO Security.User_Roles (id, user_id, role_id) VALUES (3, 2, 1);
INSERT INTO Security.User_Roles (id, user_id, role_id) VALUES (4, 3, 2);
INSERT INTO Security.User_Roles (id, user_id, role_id) VALUES (5, 4, 4);
INSERT INTO Security.User_Roles (id, user_id, role_id) VALUES (6, 5, 4);
SET IDENTITY_INSERT Security.User_Roles OFF;

SET IDENTITY_INSERT Consulting.Categories ON;
INSERT INTO Consulting.Categories (id, name, description) VALUES (1, 'Networking', 'Any general networking project');
INSERT INTO Consulting.Categories (id, name, description) VALUES (2, 'Database', 'Database-related projects');
INSERT INTO Consulting.Categories (id, name, description) VALUES (3, 'Software (Development)', 'Software development-related projects');
INSERT INTO Consulting.Categories (id, name, description) VALUES (4, 'Software (Integrations)', 'Integrating existing software into a user''s infrastructure');
INSERT INTO Consulting.Categories (id, name, description) VALUES (5, 'Software (Support)', 'Supporting current software offerings');
SET IDENTITY_INSERT Consulting.Categories OFF;

SET IDENTITY_INSERT Consulting.Project_Status ON;
INSERT INTO Consulting.Project_Status (id, status, status_display) VALUES (1, 'open', 'Open');
INSERT INTO Consulting.Project_Status (id, status, status_display) VALUES (2, 'completed', 'Completed');
INSERT INTO Consulting.Project_Status (id, status, status_display) VALUES (3, 'pending', 'Pending');
INSERT INTO Consulting.Project_Status (id, status, status_display) VALUES (4, 'hold', 'On Hold');
INSERT INTO Consulting.Project_Status (id, status, status_display) VALUES (5, 'cancelled', 'Cancelled');
SET IDENTITY_INSERT Consulting.Project_Status OFF;

SET IDENTITY_INSERT Consulting.Project_Task_Status ON;
INSERT INTO Consulting.Project_Task_Status (id, status, status_display) VALUES (1, 'open', 'Open');
INSERT INTO Consulting.Project_Task_Status (id, status, status_display) VALUES (2, 'completed', 'Completed');
INSERT INTO Consulting.Project_Task_Status (id, status, status_display) VALUES (3, 'pending', 'Pending');
INSERT INTO Consulting.Project_Task_Status (id, status, status_display) VALUES (4, 'hold', 'On Hold');
INSERT INTO Consulting.Project_Task_Status (id, status, status_display) VALUES (5, 'cancelled', 'Cancelled');
SET IDENTITY_INSERT Consulting.Project_Task_Status OFF;

SET IDENTITY_INSERT Consulting.Project_Task_Time ON;
INSERT INTO Consulting.Project_Task_Time (id, task_id, time, person_completed, completed_timestamp, notes) VALUES (7, 1, 2.00, 1, '2018-10-15 00:30:07.907', 'Backed-up/migrated data');
INSERT INTO Consulting.Project_Task_Time (id, task_id, time, person_completed, completed_timestamp, notes) VALUES (8, 2, 9.00, 1, '2018-10-15 00:30:38.750', 'Transferred metadata/transformed to and built new schema in SQL Server');
INSERT INTO Consulting.Project_Task_Time (id, task_id, time, person_completed, completed_timestamp, notes) VALUES (14, 4, 10.00, 1, '2018-10-15 00:33:02.457', 'Transformed data, removed unnecessary columns per customer specifications');
INSERT INTO Consulting.Project_Task_Time (id, task_id, time, person_completed, completed_timestamp, notes) VALUES (15, 5, 2.00, 1, '2018-10-15 00:33:22.800', 'Migrated data to new database engine');
INSERT INTO Consulting.Project_Task_Time (id, task_id, time, person_completed, completed_timestamp, notes) VALUES (16, 4, 4.00, 1, '2018-10-15 01:15:39.553', 'Did some more data cleanup');
INSERT INTO Consulting.Project_Task_Time (id, task_id, time, person_completed, completed_timestamp, notes) VALUES (17, 2, 3.00, 1, '2018-10-15 01:30:45.317', 'Cleaned up schema');
INSERT INTO Consulting.Project_Task_Time (id, task_id, time, person_completed, completed_timestamp, notes) VALUES (18, 5, 3.00, 1, '2018-10-15 01:31:05.207', 'Finalizing project');
SET IDENTITY_INSERT Consulting.Project_Task_Time OFF;

SET IDENTITY_INSERT Consulting.Project_Tasks ON;
INSERT INTO Consulting.Project_Tasks (id, project_id, task, due_date, assigned_user, estimated_hours, notes, status_id) VALUES (1, 1, 'Backup MySQL database', '2018-10-17', 1, 2.00, 'Run Mysqldump and copy result to offsite storage', 2);
INSERT INTO Consulting.Project_Tasks (id, project_id, task, due_date, assigned_user, estimated_hours, notes, status_id) VALUES (2, 1, 'Convert/transform schema in new location', '2018-10-24', 1, 12.00, 'Extract metadata and create necessary schemas, dbs and tables', 2);
INSERT INTO Consulting.Project_Tasks (id, project_id, task, due_date, assigned_user, estimated_hours, notes, status_id) VALUES (4, 1, 'Transform data', '2018-10-31', 1, 8.00, 'Remove any old/unnecessary fields (get this from primary contact), and transform data that doesn''t fit into new datatypes', 2);
INSERT INTO Consulting.Project_Tasks (id, project_id, task, due_date, assigned_user, estimated_hours, notes, status_id) VALUES (5, 1, 'Import data into new server', '2018-11-06', 1, 8.00, 'Shouldn''t take the full time allotted, unless there are any errors that need to be resolved', 2);
SET IDENTITY_INSERT Consulting.Project_Tasks OFF;

SET IDENTITY_INSERT Consulting.Projects ON;
INSERT INTO Consulting.Projects (id, business_id, title, description, category_id, start_date, due_date, primary_contact, team_assigned, status_id) VALUES (1, 1, 'Database Migration', 'Migration of corporate database from MySQL to MSSQL', 2, '2018-10-15', '2018-11-14', 1, 6, 1);
SET IDENTITY_INSERT Consulting.Projects OFF;

SET IDENTITY_INSERT Consulting.Team_Members ON;
INSERT INTO Consulting.Team_Members (id, team_id, user_id, team_position) VALUES (1, 6, 1, 'Database Owner');
INSERT INTO Consulting.Team_Members (id, team_id, user_id, team_position) VALUES (2, 6, 2, 'Team Lead');
INSERT INTO Consulting.Team_Members (id, team_id, user_id, team_position) VALUES (3, 2, 2, 'Network Admin');
INSERT INTO Consulting.Team_Members (id, team_id, user_id, team_position) VALUES (4, 2, 3, 'Jr Network Admin');
SET IDENTITY_INSERT Consulting.Team_Members OFF;

SET IDENTITY_INSERT Consulting.Teams ON;
INSERT INTO Consulting.Teams (id, name, description, specialization) VALUES (1, 'Integrations (Green)', 'Handles software integrations for Green branch', 'Software Integrations');
INSERT INTO Consulting.Teams (id, name, description, specialization) VALUES (2, 'Network Admins (Green)', 'Handles networking architecture, configuration and troubleshooting for Green branch clients', 'Networking');
INSERT INTO Consulting.Teams (id, name, description, specialization) VALUES (3, 'Developers (Green)', 'Software developers for Green branch', 'Development');
INSERT INTO Consulting.Teams (id, name, description, specialization) VALUES (4, 'Help Desk (Green)', 'Help desk technicians for Green branch', 'Help Desk');
INSERT INTO Consulting.Teams (id, name, description, specialization) VALUES (5, 'Analysts (Green)', 'Business analysts for Green branch', 'Business Systems');
INSERT INTO Consulting.Teams (id, name, description, specialization) VALUES (6, 'Database Admins (Green)', 'Database administrators for Green branch', 'Database');
SET IDENTITY_INSERT Consulting.Teams OFF;

SET IDENTITY_INSERT Accounts.Business_Contacts ON;
INSERT INTO Accounts.Business_Contacts (id, business_id, contact_id, location_id, primary_contact, notes) VALUES (1, 1, 1, 1, 1, null);
INSERT INTO Accounts.Business_Contacts (id, business_id, contact_id, location_id, primary_contact, notes) VALUES (2, 5, 4, 5, 1, null);
INSERT INTO Accounts.Business_Contacts (id, business_id, contact_id, location_id, primary_contact, notes) VALUES (3, 5, 3, 5, 0, null);
INSERT INTO Accounts.Business_Contacts (id, business_id, contact_id, location_id, primary_contact, notes) VALUES (4, 6, 6, 4, 1, null);
INSERT INTO Accounts.Business_Contacts (id, business_id, contact_id, location_id, primary_contact, notes) VALUES (5, 6, 7, 4, 0, null);
INSERT INTO Accounts.Business_Contacts (id, business_id, contact_id, location_id, primary_contact, notes) VALUES (6, 3, 5, null, 1, null);
SET IDENTITY_INSERT Accounts.Business_Contacts OFF;

SET IDENTITY_INSERT Accounts.Business_Locations ON;
INSERT INTO Accounts.Business_Locations (id, business_id, name, address1, address2, city, state, zip, country, phone, notes) VALUES (1, 1, 'Corporate Headquarters', '4120 Freidrich Ln', null, 'Austin', 'Texas', '78744', 'United States', '555-342-6324', 'Open 7 days a week');
INSERT INTO Accounts.Business_Locations (id, business_id, name, address1, address2, city, state, zip, country, phone, notes) VALUES (2, 2, 'Main Location', '544 Mateo Street', null, 'Philadelphia', 'Pennsylvania', '19135', 'United States', '555-483-2854', 'Owners are a bit unpredictable, may attack you');
INSERT INTO Accounts.Business_Locations (id, business_id, name, address1, address2, city, state, zip, country, phone, notes) VALUES (3, 4, 'Main Location', '199 Lafayette Street', null, 'New York City', 'New York', '10012', 'United States', '555-384-8684', 'Ask for Gunther');
INSERT INTO Accounts.Business_Locations (id, business_id, name, address1, address2, city, state, zip, country, phone, notes) VALUES (4, 6, 'Scranton Branch', '1725 Slough Avenue', null, 'Scranton', 'Pennsylvania', '18503', 'United States', '555-964-1342', 'Hank, the security guard will ask front desk to let you in');
INSERT INTO Accounts.Business_Locations (id, business_id, name, address1, address2, city, state, zip, country, phone, notes) VALUES (5, 5, 'Global Headquarters', '123 Clarendon Rd', null, 'Southampton', null, 'SO16 4GD', 'United Kingdom', '+1555-234-5321', 'If you can''t get in, try turning your badge off and on again');
SET IDENTITY_INSERT Accounts.Business_Locations OFF;

SET IDENTITY_INSERT Accounts.Businesses ON;
INSERT INTO Accounts.Businesses (id, name, description) VALUES (1, 'Initech', 'Company specializing in banking software');
INSERT INTO Accounts.Businesses (id, name, description) VALUES (2, 'Paddy''s Pub', 'A lovely little spot in the heart of Philadelphia');
INSERT INTO Accounts.Businesses (id, name, description) VALUES (3, 'Bluth Company', 'Real estate development company');
INSERT INTO Accounts.Businesses (id, name, description) VALUES (4, 'Central Perk', 'NYC coffee shop');
INSERT INTO Accounts.Businesses (id, name, description) VALUES (5, 'Reynholm Industries', 'Incredibly successful British corporation');
INSERT INTO Accounts.Businesses (id, name, description) VALUES (6, 'Dunder Mifflin', 'The people person''s paper people');
INSERT INTO Accounts.Businesses (id, name, description) VALUES (7, 'Pied Piper, Inc.', 'Compression company');
SET IDENTITY_INSERT Accounts.Businesses OFF;

SET IDENTITY_INSERT Accounts.Contacts ON;
INSERT INTO Accounts.Contacts (id, salutation, first_name, middle_name, last_name, title, email, phone_cell, phone_office, preferred_method, notes, business_id) VALUES (1, 'Mr', 'Bill', null, 'Lumbergh', 'Manager', 'blumberg@initech.co', '555-324-1643', '555-848-8483 (x342)', 'Office Phone', 'Make sure you fill out your TPS reports', null);
INSERT INTO Accounts.Contacts (id, salutation, first_name, middle_name, last_name, title, email, phone_cell, phone_office, preferred_method, notes, business_id) VALUES (2, 'Mr', 'Richard', null, 'Hendricks', 'CEO', 'rhendricks@piedpiper.com', '555-843-8382', '555-990-8493 (x443)', 'Cell Phone', 'Looking for investors', null);
INSERT INTO Accounts.Contacts (id, salutation, first_name, middle_name, last_name, title, email, phone_cell, phone_office, preferred_method, notes, business_id) VALUES (3, 'Ms', 'Jen ', null, 'Barber', 'IT Manager', 'jbarber@reynhom.co.uk', '555-483-7784', '555-445-2543', 'Cell Phone', 'Very knowledgeable about Internet Things', null);
INSERT INTO Accounts.Contacts (id, salutation, first_name, middle_name, last_name, title, email, phone_cell, phone_office, preferred_method, notes, business_id) VALUES (4, 'Mr', 'Douglas', null, 'Reynholm', 'CEO', 'dreynholm@reynholm.co.uk', '555-555-2433', '555-884-4953', 'Cell Phone', 'Can be very handsy, be careful', null);
INSERT INTO Accounts.Contacts (id, salutation, first_name, middle_name, last_name, title, email, phone_cell, phone_office, preferred_method, notes, business_id) VALUES (5, 'Mr', 'George', null, 'Bluth', 'Founder', 'gbluthsr@bluth.com', '555-939-4832', null, 'Office Phone', 'Last seen in the walls', null);
INSERT INTO Accounts.Contacts (id, salutation, first_name, middle_name, last_name, title, email, phone_cell, phone_office, preferred_method, notes, business_id) VALUES (6, 'Mr', 'Michael', null, 'Scott', 'Regional Manager', 'mscott@dundermifflin.com', '555-493-6852', '555-889-4732 (x802)', 'Office Phone', 'Talks a lot', null);
INSERT INTO Accounts.Contacts (id, salutation, first_name, middle_name, last_name, title, email, phone_cell, phone_office, preferred_method, notes, business_id) VALUES (7, 'Mr', 'Dwight', null, 'Schrute', 'Assistant (to the) Regional Manager', 'dschrute@dundermifflin.com', '555-334-2453', null, 'Cell Phone', 'Just the facts, ma''am', null);
SET IDENTITY_INSERT Accounts.Contacts OFF;


EXEC sp_msforeachtable 'ALTER TABLE ? WITH CHECK CHECK CONSTRAINT all';


