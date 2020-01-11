# file: features/list.feature

Feature: List appointments for a day/week/month

  Background:
    Given the wall time is "2020-01-02 12:34"
    And appointment "X" owned by "Mr.EA" starts at "2020-01-02 15:00"
    And appointment "Y" owned by "Mr.Who" starts at "2020-01-08 15:00"
    And appointment "Z" owned by "Mr.EA" starts at "2020-02-01 15:00"

  Scenario: List for day
    When I list appointments for "Mr.EA" for "day" period
    Then appointment "X" is listed
    And no other appointments are listed

  Scenario: List for week
    When I list appointments for "Mr.Who" for "week" period
    Then appointment "Y" is listed
    And no other appointments are listed

  Scenario: List for month
    When I list appointments for "Mr.EA" for "month" period
    Then appointment "X" is listed
    And appointment "Z" is listed
    And no other appointments are listed
