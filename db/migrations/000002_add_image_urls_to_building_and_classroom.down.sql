-- Remove the image_urls column from the building table
ALTER TABLE "building" DROP COLUMN "image_urls";

-- Remove the image_urls column from the classroom table
ALTER TABLE "classroom" DROP COLUMN "image_urls";
