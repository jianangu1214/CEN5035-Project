describe('TravelLog smoke', () => {
  it('logs in and creates a hotel stay', () => {
    const email = Cypress.env('TEST_EMAIL')
    const password = Cypress.env('TEST_PASSWORD')

    cy.visit('/login')
    cy.get('[data-testid="login-email"]').clear().type(email)
    cy.get('[data-testid="login-password"]').clear().type(password, { log: false })
    cy.get('[data-testid="login-submit"]').click()

    cy.location('pathname').should('eq', '/hotels')

    cy.visit('/hotels/new')
    cy.get('#hotel_name').type('Cypress Test Inn')
    cy.get('#city').type('Orlando')
    cy.get('#country').type('USA')
    cy.get('#check_in').type('2026-04-01')
    cy.get('#check_out').type('2026-04-03')
    cy.get('#price').clear().type('199.50')
    cy.get('[data-testid="hotel-submit"]').click()

    cy.location('pathname').should('eq', '/hotels')
    cy.contains('Cypress Test Inn')
  })
})
