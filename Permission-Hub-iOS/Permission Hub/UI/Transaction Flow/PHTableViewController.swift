//
//  TransactionOverviewViewController.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

enum TransactionVerificationStatus {
    case verified, unverified

    var color: UIColor {
        switch self {
        case .verified:
            return PHColors.topaz
        case .unverified:
            return PHColors.red
        }
    }
}

enum PHTableViewViewCellType {
    case notification(
        type: TransactionNotificationType,
        text: String)
    case warning(text: String)
    case description(text: String)
    case transactionItem(item: TransactionItem)
    case form
}

final class PHTableViewController: UIViewController {

    // MARK: - Private properties

    private lazy var tableView: UITableView = {

        let tableView = UITableView()
        tableView.translatesAutoresizingMaskIntoConstraints = false

        tableView.rowHeight = UITableView.automaticDimension
        tableView.estimatedRowHeight = 66
        tableView.separatorStyle = .none
        tableView.allowsMultipleSelection = true

        tableView.dataSource = self
        tableView.delegate = self

        tableView.register(
            TransactionTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionTableViewCell.self))

        tableView.register(
            TransactionDescriptionTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionDescriptionTableViewCell.self))

        tableView.register(
            TransactionNotificationTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionNotificationTableViewCell.self))

        tableView.register(
            FormTextInputCell.self,
            forCellReuseIdentifier: String(describing: FormTextInputCell.self))

        return tableView
    }()

    private lazy var bottomStackView: UIStackView = {

        let stackView = UIStackView()
        stackView.translatesAutoresizingMaskIntoConstraints = false

        stackView.axis = .horizontal
        stackView.alignment = .center
        stackView.distribution = .equalCentering
        stackView.spacing = 20

        return stackView
    }()

    private var items: [PHTableViewViewCellType]

    // MARK: - Initialization

    init(items: [PHTableViewViewCellType]) {
        self.items = items

        super.init(nibName: nil, bundle: nil)
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        view.backgroundColor = PHColors.lightGray

        view.addSubview(tableView)
        view.addSubview(bottomStackView)

        tableView.topAnchor.constraint(equalTo: view.topAnchor).isActive = true
        tableView.leftAnchor.constraint(equalTo: view.leftAnchor).isActive = true
        tableView.rightAnchor.constraint(equalTo: view.rightAnchor).isActive = true
        bottomStackView.heightAnchor.constraint(equalToConstant: 75).isActive = true
        bottomStackView.topAnchor.constraint(equalTo: tableView.bottomAnchor).isActive = true
        bottomStackView.centerXAnchor.constraint(equalTo: view.centerXAnchor).isActive = true
        bottomStackView.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor).isActive = true

        for i in 0..<2 {

            let button = UIButton()
            button.translatesAutoresizingMaskIntoConstraints = false

            button.widthAnchor.constraint(equalToConstant: 92).isActive = true
            button.heightAnchor.constraint(equalToConstant: 34).isActive = true

            button.titleLabel?.font = PHFonts.regular()

            let color = i == 0 ? PHColors.greyishBrown : PHColors.topaz
            button.setTitleColor(color, for: .normal)
            button.backgroundColor = .white

            button.layer.cornerRadius = 5
            button.layer.borderWidth = 1
            button.layer.borderColor = color.cgColor

            let title = i == 0 ? "Decline" : "Continue"
            button.setTitle(title, for: .normal)
            
            bottomStackView.addArrangedSubview(button)
        }
    }
}

extension PHTableViewController: UITableViewDataSource {

    func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return items.count
    }

    func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {

        switch items[indexPath.row] {
        case .notification(let type, let text):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionNotificationTableViewCell.self),
                for: indexPath) as! TransactionNotificationTableViewCell

            cell.configure(withType: type, andText: text)

            return cell

        case .description(let text):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionDescriptionTableViewCell.self),
                for: indexPath) as! TransactionDescriptionTableViewCell

            cell.configure(withText: text)

            return cell

        case .transactionItem(let item):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionTableViewCell.self),
                for: indexPath) as! TransactionTableViewCell

            let viewModel = TransactionTableViewCellViewModel(
                image: UIImage(),
                title: item.item,
                subtitle: (item.fields as NSArray).componentsJoined(by: ", "))
            cell.configure(withViewModel: viewModel)

            return cell

        case .warning(let text):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionTableViewCell.self),
                for: indexPath) as! TransactionTableViewCell

            let viewModel = TransactionTableViewCellViewModel(
                image: UIImage(),
                title: text,
                subtitle: "")
            cell.configure(withViewModel: viewModel)

            return cell

        case .form:

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: FormTextInputCell.self),
                for: indexPath) as! FormTextInputCell

            cell.configure(withPlaceholder: "text") { text in
                print(text)
            }

            return cell
        }
    }
}

extension PHTableViewController: UITableViewDelegate {

    func tableView(_ tableView: UITableView, didSelectRowAt indexPath: IndexPath) {

        switch items[indexPath.row] {
        case .notification:
            let items: [PHTableViewViewCellType] = [
                .warning(text: "Basic informations are not needed."),
                .warning(text: "Health records are not mandatory.")
            ]
            let viewController = PHTableViewController(items: items)
            navigationController?.pushViewController(viewController, animated: true)

        case .transactionItem:
            tableView.cellForRow(at: indexPath)?.isSelected = true

        default:
            break
        }
    }

    func tableView(_ tableView: UITableView, didDeselectRowAt indexPath: IndexPath) {

        switch items[indexPath.row] {
        case .transactionItem:
            tableView.cellForRow(at: indexPath)?.isSelected = false

        default:
            break
        }
    }
}
