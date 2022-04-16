using JuMP
using GLPK
using Dates

instances = readdir("./instancias", join=true)

struct Instance
    n :: Int # Vertices count
    e :: Int # Edges count, complete graph
    distances :: Vector{Vector{Int}} # Distances matrix of the graph
    demands :: Vector{Int} # Demand of each vertex
    limits :: Vector{Int}  # Load limit of each vertex
end

function read_instance(filepath :: String)
    n = 0 # Vertices count
    e = 0
    distances = Vector{Vector{Int}}()
    demands = Vector{Int}()
    limits = Vector{Int}()
    
    open(filepath, "r") do f
        n = parse(Int, readline(f))
        e = floor(n*(n-1)/2)
        for i = 1:n
            line = split(strip(readline(f)), " ")
            distance_line = Vector{Int}()
            for dist_str in line
                if dist_str == ""
                    continue
                end
               push!(distance_line, parse(Int, dist_str)) 
            end
            push!(distances, distance_line)
            
        end
        demands_line = split(strip(readline(f)), " ")
            limits_line = split(strip(readline(f)), " ")

            for j = 1:n
                if demands_line[j] != "" 
                    push!(demands, parse(Int, demands_line[j]))
                end

                if limits_line[j] != ""
                    push!(limits, parse(Int, limits_line[j]))
                end
            end
    end

    return Instance(n, e, distances, demands, limits)
end

function run_all_instances()
    output_results = "instance_name,n,linnear_bkv,linnear_opt,metaheurisitc_bvk,metaheurisitc_opt"
    for instance_filepath in instances
        instance_display_name = replace(instance_filepath, "./instancias\\" => "")
        if instance_display_name == "resultados.dat"
            continue
        end
        println("Running instance ", instance_display_name)
        instance = read_instance(instance_filepath)
        lps = solve_linnear_programming(instance)
        mhs = solve_metaheuristic(instance)
        output_results = string(output_results, "\n", replace(instance_display_name, ".dat" => ""), ",", instance.n, ",", lps[1], ",", lps[2], ",", mhs[1], ",", mhs[2])
        break
    end

    open(replace(string("results_", now(), ".csv"), ":" => ""), "w") do f 
        write(f, output_results)
    end
end

function solve_linnear_programming(instance :: Instance) 
    model = Model(GLPK.Optimizer)
    
    CX = sum(Int,instance.demands)
    println(typeof(CX))
    @variable(model,traveled[1:instance.n,1:instance.n],Bin)
    @variable(model,carga[1:instance.n],Int)

    @constraint(model,carga[1] <= CX - sum(Int,traveled[1,i]*instance.demands[i] for i in 1:instance.n))
    @constraint(model,CX <=  sum(Int,traveled[1,i]*instance.limits[i] for i in 1:instance.n))
    for i in 2: instance.n
    @constraint(model,carga[i] <= sum(Int,carga[x]*traveled[x,i] for x in 1:instance.n) - sum(Int,traveled[i,x]*instance.demands[i] for x in 1:instance.n))
    @constraint(model,sum(Int,carga[x]*traveled[x,i] for x in 1:instance.n) <= sum(Int,traveled[i,x]*instance.demands[i] for x in 1:instance.n))
    end

    @constraint(model, sum(Int,traveled[i,i] for i in 1:instance.n) <= 0)
    for i in 1: instance.n
    @constraint(model, sum(Int,traveled[i,j] for j in 1:instance.n) <= 1)
    @constraint(model, sum(Int,traveled[j,i] for j in 1:instance.n) <= 1)
    end

    @objective(model,Min,sum(Int,instance.distances[i,j]*traveled[i,j] for i in 1:instance.n for j in 1:instance.n))

    optimize!(model)

    @show objective_value(model)
    return [0, true]
end

function solve_metaheuristic(instance :: Instance) 
    return [0, true]
end

run_all_instances()