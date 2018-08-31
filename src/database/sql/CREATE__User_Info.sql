SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[User_Info](
	[id] [uniqueidentifier] NOT NULL,
	[user_id] [uniqueidentifier] NOT NULL,
	[preferred_name] [varchar](50) NULL,
	[gender] [varchar](10) NULL,
	[date_of_birth] [date] NULL,
	[email] [varchar](100) NOT NULL,
	[phone_primary] [varchar](20) NULL,
	[phone_secondary] [varchar](20) NULL,
	[address1] [varchar](50) NULL,
	[address2] [varchar](50) NULL,
	[city] [varchar](50) NULL,
	[state] [char](2) NULL,
	[zip] [varchar](10) NULL,
	[bio] [text] NULL
) ON [PRIMARY] TEXTIMAGE_ON [PRIMARY]
GO
ALTER TABLE [dbo].[User_Info] ADD  CONSTRAINT [PK_User_Info] PRIMARY KEY CLUSTERED 
(
	[id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, SORT_IN_TEMPDB = OFF, IGNORE_DUP_KEY = OFF, ONLINE = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
GO
ALTER TABLE [dbo].[User_Info] ADD  CONSTRAINT [DF_User_Info_id]  DEFAULT (newid()) FOR [id]
GO
ALTER TABLE [dbo].[User_Info]  WITH CHECK ADD  CONSTRAINT [FK_User_Info_Users] FOREIGN KEY([user_id])
REFERENCES [dbo].[Users] ([id])
GO
ALTER TABLE [dbo].[User_Info] CHECK CONSTRAINT [FK_User_Info_Users]
GO
