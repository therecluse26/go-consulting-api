SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Employee_Info](
	[id] [uniqueidentifier] NOT NULL,
	[user_id] [uniqueidentifier] NOT NULL,
	[title] [varchar](50) NOT NULL,
	[department] [varchar](50) NULL,
	[manager_id] [uniqueidentifier] NULL,
	[start_date] [date] NOT NULL
) ON [PRIMARY]
GO
ALTER TABLE [dbo].[Employee_Info] ADD  CONSTRAINT [PK_Employee_Info] PRIMARY KEY CLUSTERED 
(
	[id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, SORT_IN_TEMPDB = OFF, IGNORE_DUP_KEY = OFF, ONLINE = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
GO
ALTER TABLE [dbo].[Employee_Info] ADD  CONSTRAINT [DF_Employee_Info_id]  DEFAULT (newid()) FOR [id]
GO
ALTER TABLE [dbo].[Employee_Info] ADD  CONSTRAINT [DF_Employee_Info_start_date]  DEFAULT (CONVERT([date],getdate())) FOR [start_date]
GO
