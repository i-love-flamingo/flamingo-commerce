# 14. June 2018

* Price Fields in Cartitems and Carttotals have been changed:
  * Cartitem:
    * Deleted (Dont use anymore): Price / DiscountAmount / PriceInclTax
    * Now Existing: SinglePrice / SinglePriceInclTax / RowTotal / TaxAmount/ RowTotalInclTax / TotalDiscountAmount / ItemRelatedDiscountAmount / NonItemRelatedDiscountAmount / RowTotalWithItemRelatedDiscount / RowTotalWithItemRelatedDiscountInclTax / RowTotalWithDiscountInclTax
    
  * Carttotal:
    * Deleted: DiscountAmount
    * Now Existing: SubTotal / SubTotalInclTax / SubTotalInclTax /SubTotalWithDiscounts / SubTotalWithDiscountsAndTax / TotalDiscountAmount / TotalNonItemRelatedDiscountAmount 
