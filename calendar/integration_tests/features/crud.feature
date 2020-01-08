# file: features/crud.feature

Feature: CRUD operations on appointments

  I should be able to create, modify, delete and list appointments

  Scenario: New appointment registered
    When I send Create request
    Then The response is OK
    And the response contains the generated id

  Scenario: Appointment is shown by its UUID
    Given some appointment is registered
    When I send GetById request for given id
    Then I receive the valid properties

  Scenario: Delete appointment
    Given some appointment is registered
    When I send Delete request for given id
    And I send GetById request for given id
    Then it returns ErrNoSuchId

  Scenario: Update an existing appointment
    Given some appointment is registered
    When provide another owner "Mr.Who"
    And I send Update request for given id
    Then appointment's owner is "Mr.Who"
    But other properies are not
