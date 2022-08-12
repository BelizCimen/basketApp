# The functions of this service are as follows;

-Add Products: 
  Products can be added to system
	
-Add Products: 
  Users can be added to system
	
-List Products: 
  Users should be able to list all products.
	
-Add To Cart:
  Users can add their products to the cart 
	
-Show Cart:
  Users can list the products they have added to their cart.
	
-Delete Cart Item:
  Users can remove the products added from the cart.
	
-Complete Order:
  Users can create an order with their cart and get total price with discount and get order details.
  
  Some business rules
Products always have price and VAT (Value Added Tax, or KDV). VAT might be different for different products. Typical VAT percentage is %1, %8 and %18.
There might be discount in following situations: a. Every fourth order whose total is more than given amount may have discount depending on products. Products whose VAT is %1 don't have any discount but products whose VAT is %8 and %18 have discount of %10 and %15 respectively. b. If there are more than 3 items of the same product, then fourth and subsequent ones would have %8 off. c. If the customer made purchase which is more than given amount in a month then all subsequent purchases should have %10 off. d. Only one discount can be applied at a time. Only the highest discount should be applied.
