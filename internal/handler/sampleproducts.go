package handler

import (
	"goRent/internal/model"
	"time"
)

var (
	product1 = model.ElasticSearchProductSample{
		ID:          1,
		OwnerID:     1,
		Brand:       "Nike",
		Title:       "Nike Super Air Jordan Shoes",
		Rating:      0,
		Description: "super cheap nike super air jordan shoes for rent",
		Price:       12.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Images:      []string{"https://wooteam-productslist.s3-ap-southeast-1.amazonaws.com/product_list/images/1_1.jpeg"},
	}

	product2 = model.ElasticSearchProductSample{
		ID:          2,
		OwnerID:     1,
		Brand:       "Nike",
		Title:       "Nike Vaporfly",
		Rating:      0,
		Description: "Super Limited Edition Nike shoe. Weaar once and these shoes will make you take off. Experience the all new air jordans today!",
		Price:       29.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Images:      []string{"https://wooteam-productslist.s3-ap-southeast-1.amazonaws.com/product_list/images/2_1.jpeg"},
	}

	product3 = model.ElasticSearchProductSample{
		ID:          3,
		OwnerID:     1,
		Brand:       "Yonex",
		Title:       "Yonex Badminton Racket",
		Rating:      0,
		Description: "Want to play badminton like a pro but dont have a pro racket? Or even if you just want to play leisurely over a weekend, you can now get to use a professional badminton racket! At just $7.49 a day (a fraction of the cost you pay for a new one), you can rent a professional racket from us!",
		Price:       7.49,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Images:      []string{"https://wooteam-productslist.s3-ap-southeast-1.amazonaws.com/product_list/images/3_1.jpeg"},
	}

	product4 = model.ElasticSearchProductSample{
		ID:          4,
		OwnerID:     1,
		Brand:       "Seba",
		Title:       "Inline Roller Skates",
		Rating:      0,
		Description: "Renting out my beloved pair of skates as I no longer use them because of a new knee injury but dont wish to sell them because it carries sentimental values.",
		Price:       9.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Images:      []string{"https://wooteam-productslist.s3-ap-southeast-1.amazonaws.com/product_list/images/4_1.jpeg"},
	}

	product5 = model.ElasticSearchProductSample{
		ID:          5,
		OwnerID:     1,
		Brand:       "For Dummies",
		Title:       "GO FOR DUMMIES",
		Rating:      0,
		Description: "If youve been struggling with picking up GO as a programming language, this book will definitely change the way you think of go. Contains 101 secrets used by top professionals in the industry.",
		Price:       2.94,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Images:      []string{"https://wooteam-productslist.s3-ap-southeast-1.amazonaws.com/product_list/images/5_1.jpeg"},
	}

	product6 = model.ElasticSearchProductSample{
		ID:          6,
		OwnerID:     1,
		Brand:       "Crumpler",
		Title:       "Camera Bag",
		Rating:      0,
		Description: "A good photographer invest in their camera bodies and lensese. An Excellent one invests in protection for their gear. The crumpler 6 divisoin camera bag is bound to offer the best protection even for your top range gears. Comes with extra padding if you need them!",
		Price:       3.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Images:      []string{"https://wooteam-productslist.s3-ap-southeast-1.amazonaws.com/product_list/images/6_1.jpeg"},
	}

	product7 = model.ElasticSearchProductSample{
		ID:          7,
		OwnerID:     1,
		Brand:       "Casio",
		Title:       "Graphical calculator",
		Rating:      0,
		Description: "A calculator that calculates so fast your eyes won't even be able to see it.",
		Price:       1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Images:      []string{"https://wooteam-productslist.s3-ap-southeast-1.amazonaws.com/product_list/images/7_1.jpeg"},
	}

	product8 = model.ElasticSearchProductSample{
		ID:          8,
		OwnerID:     1,
		Brand:       "HP",
		Title:       "Monochrome Printer",
		Rating:      0,
		Description: "This printer is not just your run-of-the-mill printer. It is a monochrome printer. That means it doesnt just print in black and white. It also doesnt scream at you when you run out of red, green even when you do not use those colors",
		Price:       9.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Images:      []string{"https://wooteam-productslist.s3-ap-southeast-1.amazonaws.com/product_list/images/8_1.jpeg"},
	}

	product9 = model.ElasticSearchProductSample{
		ID:          9,
		OwnerID:     1,
		Brand:       "Dell",
		Title:       "Ultra wide monitor U2917W",
		Rating:      0,
		Description: "Have to work in excel and have to scroll countless of times to the left or right? Well, youre not alone. This monitor right here will be your best buddy when you need it for your work or assignment. With the horizontal length more than 2.5 times a regular 11 inch macbook pro, this monitor is bound to help you improve your efficiency!",
		Price:       4.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Images:      []string{"https://wooteam-productslist.s3-ap-southeast-1.amazonaws.com/product_list/images/9_1.jpeg"},
	}

	product10 = model.ElasticSearchProductSample{
		ID:          10,
		OwnerID:     1,
		Brand:       "HP",
		Title:       "school computer laptops",
		Rating:      0,
		Description: "All these computer laptops are availale for rent. They are sanizited frequently and definitely after they are returned to us. All 12 of them are good for you to host your own programming classes. If you need more quantity, please feel free to ask or PM.",
		Price:       129.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Images:      []string{"https://wooteam-productslist.s3-ap-southeast-1.amazonaws.com/product_list/images/10_1.jpeg"},
	}

	product11 = model.ElasticSearchProductSample{
		ID:          11,
		OwnerID:     1,
		Brand:       "PrettyBride",
		Title:       "pretty wedding dress",
		Rating:      0,
		Description: "I bought this dress during my wedding and feel that its such a waste if I just let it stay in the closet. It comes with these fancy shiny beads that help make you sparkle even more while you walk down the aisle! While I cannot fit into this wedding dress anymore, Ill be glad to see someone else in it",
		Price:       88.88,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Images:      []string{"https://wooteam-productslist.s3-ap-southeast-1.amazonaws.com/product_list/images/11_1.jpeg"},
	}

	product12 = model.ElasticSearchProductSample{
		ID:          12,
		OwnerID:     1,
		Brand:       "Disney",
		Title:       "lion king movie cd",
		Rating:      0,
		Description: "You know, one day our children wont even know what a CD player is. Rent this, take this time to educate them on what a CD players is. Besides, lion king is very good show about children. #ensureCDplayerGetsPassedDown",
		Price:       3.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Images:      []string{"https://wooteam-productslist.s3-ap-southeast-1.amazonaws.com/product_list/images/12_1.jpeg"},
	}

	product13 = model.ElasticSearchProductSample{
		ID:          13,
		OwnerID:     1,
		Brand:       "IceHockey Ltd",
		Title:       "Hockey Sticks for an Entire Team",
		Rating:      0,
		Description: "We recently had all the hockey sticks replaced for the entire team because of S.O.P that the management wanted. They are as good as new because we rarely used them. Very wasted to let them rot so we have decided to rent them.",
		Price:       19.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Images:      []string{"https://wooteam-productslist.s3-ap-southeast-1.amazonaws.com/product_list/images/13_1.jpeg"},
	}

	product14 = model.ElasticSearchProductSample{
		ID:          14,
		OwnerID:     1,
		Brand:       "Marvel",
		Title:       "Iron man old series dvd",
		Rating:      0,
		Description: "You know, one day our children wont even know what a CD player is. Rent this, take this time to educate them on what a CD players is. Besides, it will be a good chance for them to observe the differences between the CGI of the past and the ones on the cinema now. #ensureCDplayerGetsPassedDown",
		Price:       3.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Images:      []string{"https://wooteam-productslist.s3-ap-southeast-1.amazonaws.com/product_list/images/14_1.jpeg"},
	}
)
