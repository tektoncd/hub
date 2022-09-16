describe('🙃 E2E Tests for Hub 🙃', () => {
  it('Visits the Tekton Hub home page', () => {
    cy.visit('http://localhost:3000');
  });

  it('Should search for a resource, filter resources based on kind and select the resource', () => {
    cy.visit('http://localhost:3000');

    cy.get('[data-test="search"]').type('buildpacks');

    cy.get('[data-test="Pipeline"]').click();

    cy.get('[data-test="buildpacks"]').click();
  });

  it('Should search for resources, filter resources based on catalog and select the resource', () => {
    cy.visit('http://localhost:3000');

    cy.get('[data-test="search"]').type('buildpacks');

    cy.get('[data-test="tekton"]').click();

    cy.get('[data-test="Pipeline"]').click();

    cy.get('[data-test="buildpacks"]').click();
  });

  it('Should search for resources, filter resources based on category and select the resource', () => {
    cy.visit('http://localhost:3000');

    cy.get('[data-test="search"]').type('rsync');

    cy.get('[data-test="CLI"]').click();

    cy.get('[data-test="rsync"]').click();
  });
});
