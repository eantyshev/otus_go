# file: features/crud.feature

Feature: CRUD operations on appointments

  I should be able to create, modify, delete and list appointments

  Scenario: New appointment registered
    When I send Create request
    Then The response is OK
    And the response contains the generated id

  Scenario: Appointment is shown by its UUID
    Given appointment with id "UUID" registered
    When I send GetById request for id "UUID"
    Then I receive the valid properties

  Scenario: Delete appointment
    Given appointment with id "UUID" exists
    When I send Delete request for id "UUID"
    And I send GetById request for id "UUID"
    Then it returns ErrNoSuchId

  Scenario: Update an existing appointment
    Given appointment with id "UUID" exists
    When provide another owner "Mr.Who"
    And I send Update request for id "UUID"
    Then appointment's owner should be updated
    But other properies are not
