local schema = 'osm'
local srid = 4326

local pois = osm2pgsql.define_table({
    name = 'pois',
    schema = schema,

    ids = {
        type = 'any',
        id_column = 'osm_id',
        type_column = 'osm_type',
        create_index = 'primary_key'
    },

    columns = {
        { column = 'name',        type = 'text' },
        { column = 'category',    type = 'text',  not_null = true },
        { column = 'subcategory', type = 'text',  not_null = true },
        { column = 'tags',        type = 'jsonb', not_null = true },
        {
            column = 'geom',
            type = 'point',
            projection = srid,
            not_null = true
        }
    },

    indexes = {
        { column = 'geom', method = 'gist' },
        {
            column = { 'category', 'subcategory' },
            method = 'btree'
        }
    }
})

local areas = osm2pgsql.define_table({
    name = 'areas',
    schema = schema,

    ids = {
        type = 'any',
        id_column = 'osm_id',
        type_column = 'osm_type',
        create_index = 'primary_key'
    },

    columns = {
        { column = 'name',        type = 'text' },
        { column = 'category',    type = 'text',  not_null = true },
        { column = 'subcategory', type = 'text',  not_null = true },
        { column = 'tags',        type = 'jsonb', not_null = true },
        {
            column = 'geom',
            type = 'geometry',
            projection = srid,
            not_null = true
        }
    },

    indexes = {
        { column = 'geom', method = 'gist' },
        {
            column = { 'category', 'subcategory' },
            method = 'btree'
        }
    }
})

local function define_layer(name, geometry_type, columns)
    table.insert(columns, { column = 'tags', type = 'jsonb', not_null = true })
    table.insert(columns, {
        column = 'geom', type = geometry_type, projection = srid, not_null = true
    })
    return osm2pgsql.define_table({
        name = name,
        schema = schema,
        ids = {
            type = 'any', id_column = 'osm_id', type_column = 'osm_type',
            create_index = 'primary_key'
        },
        columns = columns,
        indexes = { { column = 'geom', method = 'gist' } }
    })
end

local buildings = define_layer('buildings', 'geometry', {
    { column = 'building_type', type = 'text', not_null = true },
    { column = 'name', type = 'text' },
    { column = 'house_number', type = 'text' },
    { column = 'street', type = 'text' },
    { column = 'postcode', type = 'text' }
})

local addresses = define_layer('addresses', 'point', {
    { column = 'house_number', type = 'text', not_null = true },
    { column = 'street', type = 'text' },
    { column = 'city', type = 'text' },
    { column = 'postcode', type = 'text' }
})

local roads = define_layer('roads', 'linestring', {
    { column = 'highway', type = 'text', not_null = true },
    { column = 'name', type = 'text' },
    { column = 'lanes', type = 'text' },
    { column = 'maxspeed', type = 'text' },
    { column = 'surface', type = 'text' }
})

local included_landuse = {
    residential = true,
    commercial = true,
    retail = true,
    industrial = true,
    recreation_ground = true,
    forest = true,
    meadow = true,
    grass = true
}

local included_leisure = {
    park = true,
    garden = true,
    nature_reserve = true
}

local included_natural = {
    wood = true,
    grassland = true
}

local function classify_poi(tags)
    if tags.shop then
        return 'shop', tags.shop
    end

    if tags.amenity then
        return 'amenity', tags.amenity
    end

    if tags.tourism then
        return 'tourism', tags.tourism
    end

    return nil, nil
end

local function classify_area(tags)
    if tags.leisure and included_leisure[tags.leisure] then
        return 'leisure', tags.leisure
    end

    if tags.landuse and included_landuse[tags.landuse] then
        return 'landuse', tags.landuse
    end

    if tags.natural and included_natural[tags.natural] then
        return 'natural', tags.natural
    end

    return nil, nil
end

local function insert_poi(object, geom)
    local category, subcategory = classify_poi(object.tags)

    if not category then
        return
    end

    pois:insert({
        name = object.tags.name,
        category = category,
        subcategory = subcategory,
        tags = object.tags,
        geom = geom
    })
end

local function insert_area(object, geom)
    local category, subcategory = classify_area(object.tags)

    if not category then
        return
    end

    areas:insert({
        name = object.tags.name,
        category = category,
        subcategory = subcategory,
        tags = object.tags,
        geom = geom
    })
end

local function insert_building(object, geom)
    local tags = object.tags
    if tags.building and tags.building ~= 'no' then
        buildings:insert({
            building_type = tags.building,
            name = tags.name,
            house_number = tags['addr:housenumber'],
            street = tags['addr:street'],
            postcode = tags['addr:postcode'],
            tags = tags,
            geom = geom
        })
    end
end

function osm2pgsql.process_node(object)
    local tags = object.tags
    local point = object:as_point()
    insert_poi(object, point)
    if tags['addr:housenumber'] then
        addresses:insert({
            house_number = tags['addr:housenumber'],
            street = tags['addr:street'],
            city = tags['addr:city'],
            postcode = tags['addr:postcode'],
            tags = tags,
            geom = point
        })
    end
end

function osm2pgsql.process_way(object)
    local tags = object.tags
    if tags.highway and tags.area ~= 'yes' then
        roads:insert({
            highway = tags.highway,
            name = tags.name,
            lanes = tags.lanes,
            maxspeed = tags.maxspeed,
            surface = tags.surface,
            tags = tags,
            geom = object:as_linestring()
        })
    end

    if not object.is_closed or object.tags.area == 'no' then
        return
    end

    local polygon = object:as_polygon()

    insert_poi(object, polygon:centroid())
    insert_area(object, polygon)
    insert_building(object, polygon)
end

function osm2pgsql.process_relation(object)
    if object.tags.type ~= 'multipolygon' then
        return
    end

    local multipolygon = object:as_multipolygon()

    insert_poi(object, multipolygon:centroid())
    insert_area(object, multipolygon)
    insert_building(object, multipolygon)
end
