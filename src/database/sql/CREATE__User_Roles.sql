SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[User_Roles](
	[id] [uniqueidentifier] NOT NULL,
	[user_id] [uniqueidentifier] NOT NULL,
	[role_id] [uniqueidentifier] NOT NULL
) ON [PRIMARY]
GO
ALTER TABLE [dbo].[User_Roles] ADD  CONSTRAINT [PK_User_Roles] PRIMARY KEY CLUSTERED 
(
	[id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, SORT_IN_TEMPDB = OFF, IGNORE_DUP_KEY = OFF, ONLINE = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
GO
CREATE UNIQUE NONCLUSTERED INDEX [UQ_Users_and_Roles] ON [dbo].[User_Roles]
(
	[role_id] ASC,
	[user_id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, SORT_IN_TEMPDB = OFF, IGNORE_DUP_KEY = OFF, DROP_EXISTING = OFF, ONLINE = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
GO
ALTER TABLE [dbo].[User_Roles] ADD  CONSTRAINT [DF_User_Roles_id]  DEFAULT (newid()) FOR [id]
GO
ALTER TABLE [dbo].[User_Roles]  WITH CHECK ADD  CONSTRAINT [FK_User_Roles_RoleID] FOREIGN KEY([role_id])
REFERENCES [dbo].[Roles] ([id])
ON UPDATE CASCADE
ON DELETE CASCADE
GO
ALTER TABLE [dbo].[User_Roles] CHECK CONSTRAINT [FK_User_Roles_RoleID]
GO
ALTER TABLE [dbo].[User_Roles]  WITH CHECK ADD  CONSTRAINT [FK_User_Roles_UserID] FOREIGN KEY([user_id])
REFERENCES [dbo].[Users] ([id])
ON UPDATE CASCADE
ON DELETE CASCADE
GO
ALTER TABLE [dbo].[User_Roles] CHECK CONSTRAINT [FK_User_Roles_UserID]
GO
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'role_id references id on Roles table' , @level0type=N'SCHEMA',@level0name=N'dbo', @level1type=N'TABLE',@level1name=N'User_Roles', @level2type=N'CONSTRAINT',@level2name=N'FK_User_Roles_RoleID'
GO
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'user_id references id on Users table' , @level0type=N'SCHEMA',@level0name=N'dbo', @level1type=N'TABLE',@level1name=N'User_Roles', @level2type=N'CONSTRAINT',@level2name=N'FK_User_Roles_UserID'
GO
