package seeder

var seed_20210326014100 = Seed{
	version: "20210326014100",
	content: `INSERT INTO "user" ("name", "email", "address", "password", "createdAt", "createdBy", "updatedAt", "updatedBy") values
  ('user', 'user@home24.com', 'Jakarta', '$2a$10$5.p1ONDftoudtkvcl/o30u.BMYbEzCfdGEqrOVz6fnXsRmRir0bGK', '2021-03-26 01:41:00', '0', '2021-03-26 01:41:00', '0')
  `,
}
