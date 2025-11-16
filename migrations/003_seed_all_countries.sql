-- Seed all supported countries
-- This migration adds country records for UK, Sweden, and Canada
-- US was already added in migration 002

INSERT INTO countries (code, name, data_source_name, data_source_url, data_source_description, data_source_requires_manual_download)
VALUES 
    ('UK', 'United Kingdom', 'Office for National Statistics', 
     'https://www.ons.gov.uk/peoplepopulationandcommunity/birthsdeathsandmarriages/livebirths/datasets/babynamesenglandandwalesbabynamesstatisticsboys',
     'Baby name statistics for England and Wales from the Office for National Statistics. Data is provided in separate Excel files for boys and girls, requiring manual download and Excel parsing.',
     true),
    ('SE', 'Sweden', 'Statistics Sweden (SCB)', 
     'https://www.scb.se/hitta-statistik/statistik-efter-amne/befolkning/amnesovergripande-statistik/namnstatistik/',
     'Swedish baby name statistics from Statistics Sweden (Statistiska centralbyrån). Data includes traditional Swedish names with characters like å, ä, ö. Available in Excel or CSV format.',
     true),
    ('CA', 'Canada', 'Statistics Canada', 
     'https://www150.statcan.gc.ca/n1/en/type/data',
     'Canadian baby name data by province from Statistics Canada. Data format and availability may vary by province.',
     true)
ON CONFLICT (code) DO NOTHING;

-- Update US country description to be more detailed
UPDATE countries 
SET 
    data_source_description = 'Social Security Administration baby name data covering years 1880-2024. Data is provided in CSV format without headers, with one file per year. Each record contains name, gender, and count of occurrences.',
    data_source_requires_manual_download = true
WHERE code = 'US';