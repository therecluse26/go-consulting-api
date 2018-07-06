SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Course_Sessions](
	[id] [uniqueidentifier] NOT NULL,
	[course_id] [uniqueidentifier] NOT NULL,
	[session_number] [int] NOT NULL,
	[title] [varchar](50) NULL,
	[description] [varchar](500) NULL,
	[start_datetime] [datetime] NULL
) ON [PRIMARY]
GO
ALTER TABLE [dbo].[Course_Sessions] ADD  CONSTRAINT [PK_Course_Sessions] PRIMARY KEY CLUSTERED 
(
	[id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, SORT_IN_TEMPDB = OFF, IGNORE_DUP_KEY = OFF, ONLINE = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
GO
ALTER TABLE [dbo].[Course_Sessions] ADD  CONSTRAINT [DF_Course_Sessions_id]  DEFAULT (newid()) FOR [id]
GO
ALTER TABLE [dbo].[Course_Sessions]  WITH CHECK ADD  CONSTRAINT [FK_Course_Sessions_Courses] FOREIGN KEY([course_id])
REFERENCES [dbo].[Courses] ([id])
GO
ALTER TABLE [dbo].[Course_Sessions] CHECK CONSTRAINT [FK_Course_Sessions_Courses]
GO
