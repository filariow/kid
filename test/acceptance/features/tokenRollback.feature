Feature: Rotate Token

    Scenario: Name is available
        Given Identity "sa" is created
        And   Service Account "sa" exists
        And   Secret "sa-key-1" exists
        And   Token rotation begins for identity "sa"
        And   Secret "sa-key-2" exists
        And   Token rotation completes for identity "sa"
        And   Secret "sa-key-1" does not exist
        When  Token for identity "sa" with version "1" is rolled back
        Then  Secret "sa-key-1" exists
