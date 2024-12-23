describe('User Event Registration Flow', () => {
  it('should log in, register for an event, and view event applications', () => {

    cy.visit('http://localhost:8000/login');

    cy.get('input[name="email"]').type('m@m.com');
    cy.get('input[name="password"]').type('12345678');
    cy.get('button[type="submit"]').click();

    cy.wait(1000);

    cy.reload();


    cy.get('.event-card').first().click(); 


    cy.get('button.event-register').click(); 

    cy.get('button.close-button').click();

    cy.reload();
    cy.wait(1000);


    cy.visit('http://localhost:8000/myevents');

    cy.get('.event-application').should('have.length.greaterThan', 0); 
  });
});