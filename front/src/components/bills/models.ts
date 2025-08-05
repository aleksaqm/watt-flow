export interface Pricelist {
  id: number,
  red: number,
  blue: number,
  green: number,
  tax: number,
  bill_power: number,
  date: Date,
}

export interface Bill {
  id: number,
  date: string,
  issue_date: Date,
  status: string
}

export interface SearchBill {
  id: number;
  issue_date: string;
  billing_date: string;
  pricelist: PriceListSearch;
  spent_power: number;
  price: number;
  owner: User;
  status: 'Delivered' | 'Paid'
  household: Household;
  payment_reference: string;
}

export interface PriceListSearch {
  id: number,
  valid_from: string,
  blue_zone: number,
  red_zone: number,
  green_zone: number,
  billing_power: number,
  tax: number
}

export interface User {
  id: number;
  username: string;
}

export interface Household{
  id: number,
  cadastral_number: string,
}