Feature: Rotate Token

    Scenario: Rotation begins
        Given Identity "sa" is created
        And   Service Account "sa" exists
        And   Secret "sa-key-1" exists
        When  Token rotation begins for identity "sa"
        Then  Secret "sa-key-2" exists

    Scenario: Name is available
        Given Identity "sa" is created
        And   Service Account "sa" exists
        And   Secret "sa-key-1" exists
        And   Token rotation begins for identity "sa"
        And   Secret "sa-key-2" exists
        When  Token rotation completes for identity "sa"
        Then  Secret "sa-key-1" does not exist
