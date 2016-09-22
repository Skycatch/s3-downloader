# Tasque Downloader Prototype

## Notes

This is a prototype, and has not been written to be deployed in a reusable manner.

Initial approach was to create an executable with multiple adapters - S3, Box, HTTP, etc.
That turned out to be a bad direction:
- Why pack so much complexity into an otherwise simple operation?
- Updating would affect more downstream users/applications if it served multiple purposes

Instead, the initial release will be an S3 downloader with a standard execution API.
This will allow other binaries to be created for their respective services in a robust manner.

User UI niceties will be secondary to performance.

This application will require some system level modifications in order to achieve it's maximum performance.

## Findings
### Download Speed Comparison

![screen shot 2016-09-06 at 4 06 23 pm](https://cloud.githubusercontent.com/assets/47128/18406318/422f2afa-76b0-11e6-8458-b0fbb2968f47.png)

**New Downloader Summary:**
Completed in **0.804s** for a download speed of **2.278 Gb/s** (gigabits per second)

**Node Summary:**
Fetch Keys
2016-09-09T21:08:01.472Z
2016-09-09T21:08:04.935Z
**3.463s**

Download
2016-09-09T21:08:04.999Z
2016-09-09T21:09:00.875Z
**55.876s**
