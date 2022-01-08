IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[door_schedulers]') AND type in (N'U'))
DROP TABLE [dbo].[door_schedulers]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[doorlocks]') AND type in (N'U'))
DROP TABLE [dbo].[doorlocks]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[gateways]') AND type in (N'U'))
DROP TABLE [dbo].[gateways]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[areas]') AND type in (N'U'))
DROP TABLE [dbo].[areas]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[student_schedulers]') AND type in (N'U'))
DROP TABLE [dbo].[student_schedulers]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[employee_schedulers]') AND type in (N'U'))
DROP TABLE [dbo].[employee_schedulers]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[customer_schedulers]') AND type in (N'U'))
DROP TABLE [dbo].[customer_schedulers]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[employees]') AND type in (N'U'))
DROP TABLE [dbo].[employees]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[students]') AND type in (N'U'))
DROP TABLE [dbo].[students]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[customers]') AND type in (N'U'))
DROP TABLE [dbo].[customers]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[schedulers]') AND type in (N'U'))
DROP TABLE [dbo].[schedulers]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[passwords]') AND type in (N'U'))
DROP TABLE [dbo].[passwords]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[gateway_logs]') AND type in (N'U'))
DROP TABLE [dbo].[gateway_logs]

GO
