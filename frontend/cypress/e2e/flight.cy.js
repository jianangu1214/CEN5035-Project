describe('Flights feature test', () => {
  it('logs in, opens flights page, and opens add flight form', () => {
    cy.visit('/login')

    cy.get('input[type="email"]').type('championwaw@gmail.com')
    cy.get('input[type="password"]').type('12345678')
    cy.get('form button[type="submit"]').click()

    // JWT integration + protected access
    cy.contains('Flights').click()
    cy.url().should('include', '/flights')

    // Flight list page
    cy.contains('Flights')
    cy.contains('Add flight')

    // Flight create form
    cy.contains('Add flight').click()
    cy.url().should('include', '/flights/new')
    cy.get('[name="airline"]').should('exist')
    cy.get('[name="flight_number"]').should('exist')
    cy.get('[name="from_airport"]').should('exist')
    cy.get('[name="to_airport"]').should('exist')
  })
})