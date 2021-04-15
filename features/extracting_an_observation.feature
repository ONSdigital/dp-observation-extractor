Feature: Extracting Observations and putting out as messageso

    Feature Description
    Scenario: publish a message read from a csv
        Given I have a file in a bucket with name "test_file.csv" and content:
            """
            header1,header2
            value1,value2
            """
        When I recieve a message containing the file name "s3://bucket_name/test_file.csv"
        Then these messages are sent to the output message queue:
            | Row                             |
            | ,,header1,value1,header2,value2 |