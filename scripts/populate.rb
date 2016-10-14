require 'csv'
require 'json'

# curl -X POST -d '{"name": "Twitter Inc.", "symbol": "TWTR", "exchange": "NYSE"}' http://localhost:8080/api/companies

filename = ARGV[0]
endpoint = ARGV[1]

if filename.to_s.empty? || endpoint.to_s.empty?
  puts "Usage: ruby ./populate.rb companies.csv https://companies-146403.appspot.com/api/companies"
  exit 1
end

lines = CSV.parse(IO.read(filename))

lines.each do |data|
  hash = {name: data[0], symbol: data[1], exchange: data[2]}
  cmd = "curl -s -X POST -d '#{hash.to_json}' #{endpoint}"
  puts cmd
  `#{cmd}`
end
