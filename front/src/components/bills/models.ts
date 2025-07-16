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
