This directory should contain specification for service side components in {{.ProjectName}}. The following files/directories further classify the specifications:

- main.yml: Contains the swagger spec for your REST interface with workflow extensions.
- workflows: These contain the flow DSL describing the different workflows which can be linked to different views in the swagger spec using the workflow extensions.
- templates: Contains all the templates available for post workflow processing. These templates are linked to different views using the template extensions in the swagger spec.